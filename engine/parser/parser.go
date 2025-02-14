// Package parser implements a parser for SQL statements
//
// Inspired by go/parser
package parser

import (
	"fmt"

	"github.com/kokizzu/ramsql/engine/log"
)

// The parser structure holds the parser's internal state.
type parser struct {
	tokens   []Token
	tokenLen int
	index    int

	i []Instruction
}

// NewParser returns a parser object initialized with a list of tokens
func NewParser(tokens []Token) parser {
	p := parser{
		tokens: stripSpaces(tokens),
		index:  0,
		i:      []Instruction{},
	}

	p.tokenLen = len(p.tokens)

	return p
}

func (p *parser) parse() ([]Instruction, error) {
	for p.hasNext() {
		// fmt.Printf("Token index : %d\n", p.index)

		// Found a new instruction
		if p.cur().Token == SemicolonToken {
			p.index++
			continue
		}

		// Ignore space token, not needed anymore
		if p.cur().Token == SpaceToken {
			p.index++
			continue
		}

		// Now,
		// Create a logical tree of all tokens
		// We start with first order query
		// CREATE, SELECT, INSERT, UPDATE, DELETE, TRUNCATE, DROP, EXPLAIN
		switch p.cur().Token {
		case CreateToken:
			i, err := p.parseCreate()
			if err != nil {
				return nil, err
			}
			p.i = append(p.i, *i)
			break
		case SelectToken:
			i, err := p.parseSelect()
			if err != nil {
				return nil, err
			}
			p.i = append(p.i, *i)
			break
		case InsertToken:
			i, err := p.parseInsert()
			if err != nil {
				return nil, err
			}
			p.i = append(p.i, *i)
			break
		case UpdateToken:
			i, err := p.parseUpdate()
			if err != nil {
				return nil, err
			}
			p.i = append(p.i, *i)
			break
		case DeleteToken:
			i, err := p.parseDelete()
			if err != nil {
				return nil, err
			}
			p.i = append(p.i, *i)
			break
		case TruncateToken:
			i, err := p.parseTruncate()
			if err != nil {
				return nil, err
			}
			p.i = append(p.i, *i)
			break
		case DropToken:
			log.Debug("HEY DROP HERE !\n")
			i, err := p.parseDrop()
			if err != nil {
				return nil, err
			}
			p.i = append(p.i, *i)
			break
		case ExplainToken:
			break
		case GrantToken:
			i := &Instruction{}
			i.Decls = append(i.Decls, NewDecl(Token{Token: GrantToken}))
			p.i = append(p.i, *i)
			return p.i, nil
		default:
			return nil, fmt.Errorf("Parsing error near <%s>", p.cur().Lexeme)
		}
	}

	return p.i, nil
}

func (p *parser) parseUpdate() (*Instruction, error) {
	i := &Instruction{}

	// Set UPDATE decl
	updateDecl, err := p.consumeToken(UpdateToken)
	if err != nil {
		return nil, err
	}
	i.Decls = append(i.Decls, updateDecl)

	// should be table name
	nameDecl, err := p.parseQuotedToken()
	if err != nil {
		return nil, err
	}
	updateDecl.Add(nameDecl)

	// should be SET
	setDecl, err := p.consumeToken(SetToken)
	if err != nil {
		return nil, err
	}
	updateDecl.Add(setDecl)

	// should be a list of equality
	gotClause := false
	for p.cur().Token != WhereToken {

		if !p.hasNext() && gotClause {
			break
		}

		attributeDecl, err := p.parseCondition()
		if err != nil {
			return nil, err
		}
		setDecl.Add(attributeDecl)
		p.consumeToken(CommaToken)

		// Got at least one clause
		gotClause = true
	}

	err = p.parseWhere(updateDecl)
	if err != nil {
		return nil, err
	}

	return i, nil
}

func (p *parser) parseInsert() (*Instruction, error) {
	i := &Instruction{}

	// Required: INSERT
	insertDecl, err := p.consumeToken(InsertToken)
	if err != nil {
		return nil, err
	}
	i.Decls = append(i.Decls, insertDecl)

	// Required: INTO
	intoDecl, err := p.consumeToken(IntoToken)
	if err != nil {
		return nil, err
	}
	insertDecl.Add(intoDecl)

	// Required: <TABLE-NAME>
	tableDecl, err := p.parseQuotedToken()
	if err != nil {
		return nil, err
	}
	intoDecl.Add(tableDecl)

	// Required: '('
	_, err = p.consumeToken(BracketOpeningToken)
	if err != nil {
		return nil, err
	}

	// Required: <ATTRIBUTE-NAME> [, <ATTRIBUTE-NAME>]* ')'
	for {
		decl, err := p.parseQuotedToken()
		if err != nil {
			return nil, err
		}
		tableDecl.Add(decl)

		if p.is(BracketClosingToken) {
			if _, err = p.consumeToken(BracketClosingToken); err != nil {
				return nil, err
			}

			break
		}

		_, err = p.consumeToken(CommaToken)
		if err != nil {
			return nil, err
		}
	}

	// Required: VALUES
	valuesDecl, err := p.consumeToken(ValuesToken)
	if err != nil {
		return nil, err
	}
	insertDecl.Add(valuesDecl)

	// Required: '('
	_, err = p.consumeToken(BracketOpeningToken)
	if err != nil {
		return nil, err
	}

	// Required: <ATTRIBUTE-VALUE> [, <ATTRIBUTE-VALUE>]* ')'
	for {
		decl, err := p.parseListElement()
		if err != nil {
			return nil, err
		}
		valuesDecl.Add(decl)

		if p.is(BracketClosingToken) {
			p.consumeToken(BracketClosingToken)
			break
		}

		_, err = p.consumeToken(CommaToken)
		if err != nil {
			return nil, err
		}
	}

	// Optional: RETURNING ...
	if retDecl, err := p.consumeToken(ReturningToken); err == nil {
		insertDecl.Add(retDecl)
		attrDecl, err := p.parseAttribute()
		if err != nil {
			return nil, err
		}
		retDecl.Add(attrDecl)
	}

	return i, nil
}

func (p *parser) parseType() (*Decl, error) {
	typeDecl, err := p.consumeToken(StringToken)
	if err != nil {
		return nil, err
	}

	// Maybe a complex type
	if p.is(BracketOpeningToken) {
		_, err = p.consumeToken(BracketOpeningToken)
		if err != nil {
			return nil, err
		}
		sizeDecl, err := p.consumeToken(NumberToken)
		if err != nil {
			return nil, err
		}
		typeDecl.Add(sizeDecl)
		_, err = p.consumeToken(BracketClosingToken)
		if err != nil {
			return nil, err
		}
	}

	return typeDecl, nil
}

func (p *parser) parseOrderBy(selectDecl *Decl) error {
	orderDecl, err := p.consumeToken(OrderToken)
	if err != nil {
		return err
	}
	selectDecl.Add(orderDecl)

	_, err = p.consumeToken(ByToken)
	if err != nil {
		return err
	}

	// parse attribute now
	attrDecl, err := p.parseAttribute()
	if err != nil {
		return err
	}
	orderDecl.Add(attrDecl)

	// Parse multiple ordering
	for p.cur().Token == CommaToken {
		_, err := p.consumeToken(CommaToken)
		if err != nil {
			return nil
		}

		// parse attribute now
		attrDecl, err := p.parseAttribute()
		if err != nil {
			return err
		}
		orderDecl.Add(attrDecl)
	}

	// ASC ? DESC ? nothing ?
	t := p.cur().Token
	if t == AscToken || t == DescToken {
		decl, err := p.consumeToken(AscToken, DescToken)
		if err != nil {
			return err
		}
		orderDecl.Add(decl)
	}

	return nil
}

func (p *parser) parseWhere(selectDecl *Decl) error {

	// May be WHERE  here
	// Can be ORDER BY if WHERE cause if implicit
	whereDecl, err := p.consumeToken(WhereToken)
	if err != nil {
		return err
	}
	selectDecl.Add(whereDecl)

	// Now should be a list of: Attribute and Operator and Value
	gotClause := false
	for {
		if !p.hasNext() && gotClause {
			break
		}

		if p.is(OrderToken, LimitToken, ForToken) {
			break
		}

		attributeDecl, err := p.parseCondition()
		if err != nil {
			return err
		}
		whereDecl.Add(attributeDecl)

		if p.is(AndToken, OrToken) {
			linkDecl, err := p.consumeToken(p.cur().Token)
			if err != nil {
				return err
			}
			whereDecl.Add(linkDecl)
		}

		// Got at least one clause
		gotClause = true
	}

	return nil
}

// parseBuiltinFunc looks for COUNT,MAX,MIN
func (p *parser) parseBuiltinFunc() (*Decl, error) {
	var d *Decl
	var err error

	// COUNT(attribute)
	if p.is(CountToken) {
		d, err = p.consumeToken(CountToken)
		if err != nil {
			return nil, err
		}
		// Bracket
		_, err = p.consumeToken(BracketOpeningToken)
		if err != nil {
			return nil, err
		}
		// Attribute
		attr, err := p.parseAttribute()
		if err != nil {
			return nil, err
		}
		d.Add(attr)
		// Bracket
		_, err = p.consumeToken(BracketClosingToken)
		if err != nil {
			return nil, err
		}
	}

	return d, nil
}

// parseAttribute parse an attribute of the form
// table.foo
// table.*
// "table".foo
// foo
func (p *parser) parseAttribute() (*Decl, error) {
	quoted := false
	quoteToken := DoubleQuoteToken

	if p.is(DoubleQuoteToken) || p.is(BacktickToken) {
		quoted = true
		quoteToken = p.cur().Token
		if err := p.next(); err != nil {
			return nil, err
		}
	}

	// shoud be a StringToken here
	// If there is a point after, it's a table name,
	// if not, it's the attribute
	if !p.is(StringToken, StarToken) {
		return nil, p.syntaxError()
	}
	attrDecl := NewDecl(p.cur())

	if quoted {
		// Check there is a closing quote
		if _, err := p.mustHaveNext(quoteToken); err != nil {
			log.Debug("parseAttribute: Missing closing quote")
			return nil, err
		}
	}
	// If no next token,and not quoted, then it was the atribute name
	if err := p.next(); err != nil {
		return attrDecl, nil
	}

	// Optional: '.'
	if p.is(PeriodToken) {
		_, err := p.consumeToken(PeriodToken)
		if err != nil {
			return nil, err
		}
		// if so, next must be the attribute name or a star
		tableDecl := attrDecl
		attrDecl, err = p.consumeToken(StringToken, StarToken)
		if err != nil {
			return nil, err
		}

		attrDecl.Add(tableDecl)
	}

	// Optional: AS ...
	if p.is(AsToken) {
		asDecl, err := p.consumeToken(AsToken)
		if err != nil {
			return nil, err
		}

		// Required: <ATTRIBUTE-RENAME>
		renameDecl, err := p.consumeToken(StringToken)
		if err != nil {
			return nil, err
		}

		attrDecl.Add(asDecl)
		asDecl.Add(renameDecl)
	}

	return attrDecl, nil
}

// parseQuotedToken parse a token of the form <STRING>, '<STRING>', "<STRING>", `<STRING>`
func (p *parser) parseQuotedToken() (*Decl, error) {
	quoted := false
	quoteToken := DoubleQuoteToken

	if p.is(SimpleQuoteToken) || p.is(DoubleQuoteToken) || p.is(BacktickToken) {
		quoted = true
		quoteToken = p.cur().Token
		if err := p.next(); err != nil {
			return nil, err
		}
	}

	// shoud be a StringToken or keyword token
	if !p.is(StringToken) && !p.cur().IsAWord() {
		return nil, p.syntaxError()
	}
	decl := &Decl{
		Token:  StringToken,
		Lexeme: p.cur().Lexeme,
	}

	if quoted {
		// Check there is a closing quote
		if _, err := p.mustHaveNext(quoteToken); err != nil {
			return nil, err
		}
	}

	p.next()
	return decl, nil
}

func (p *parser) parseCondition() (*Decl, error) {

	// We may have the WHERE 1 condition
	if t := p.cur(); t.Token == NumberToken && t.Lexeme == "1" {
		attributeDecl := NewDecl(t)
		p.next()
		// in case of 1 = 1
		if p.cur().Token == EqualityToken {
			t, err := p.isNext(NumberToken)
			if err == nil && t.Lexeme == "1" {
				p.consumeToken(EqualityToken)
				p.consumeToken(NumberToken)
			}
		}
		return attributeDecl, nil
	}

	// Attribute
	attributeDecl, err := p.parseAttribute()
	if err != nil {
		return nil, err
	}

	switch p.cur().Token {
	case EqualityToken, LeftDipleToken, RightDipleToken, LessOrEqualToken, GreaterOrEqualToken:
		decl, err := p.consumeToken(p.cur().Token)
		if err != nil {
			return nil, err
		}
		attributeDecl.Add(decl)
		break
	case InToken:
		inDecl, err := p.parseIn()
		if err != nil {
			return nil, err
		}
		attributeDecl.Add(inDecl)
		return attributeDecl, nil
	case IsToken:
		log.Debug("parseCondition: IsToken\n")
		decl, err := p.consumeToken(IsToken)
		if err != nil {
			return nil, err
		}
		attributeDecl.Add(decl)
		if p.cur().Token == NotToken {
			log.Debug("parseCondition: NotToken\n")
			notDecl, err := p.consumeToken(NotToken)
			if err != nil {
				return nil, err
			}
			decl.Add(notDecl)
		}
		if p.cur().Token == NullToken {
			log.Debug("parseCondition: NullToken\n")
			nullDecl, err := p.consumeToken(NullToken)
			if err != nil {
				return nil, err
			}
			decl.Add(nullDecl)
		}
		return attributeDecl, nil
	}

	// Value
	valueDecl, err := p.parseValue()
	if err != nil {
		return nil, err
	}
	attributeDecl.Add(valueDecl)
	return attributeDecl, nil
}

func (p *parser) parseIn() (*Decl, error) {
	inDecl, err := p.consumeToken(InToken)
	if err != nil {
		return nil, err
	}

	// bracket opening
	_, err = p.consumeToken(BracketOpeningToken)
	if err != nil {
		return nil, err
	}

	// list of value
	gotList := false
	for {
		v, err := p.parseValue()
		if err != nil {
			return nil, err
		}
		inDecl.Add(v)
		gotList = true

		if p.is(BracketClosingToken) {
			if gotList == false {
				return nil, fmt.Errorf("IN clause: empty list of value")
			}
			p.consumeToken(BracketClosingToken)
			break
		}

		_, err = p.consumeToken(CommaToken)
		if err != nil {
			return nil, err
		}
	}

	return inDecl, nil
}

func (p *parser) parseValue() (*Decl, error) {
	debug("parseValue")
	defer debug("~parseValue")
	quoted := false

	if p.is(SimpleQuoteToken) || p.is(DoubleQuoteToken) {
		quoted = true
		debug("value %v is quoted!", p.cur())
		_, err := p.consumeToken(SimpleQuoteToken, DoubleQuoteToken)
		if err != nil {
			return nil, err
		}
	}

	valueDecl, err := p.consumeToken(StringToken, NumberToken, TrueToken, FalseToken, DateToken, NowToken, LocalTimestampToken, NullToken)
	if err != nil {
		debug("parseValue: Wasn't expecting %v\n", p.cur())
		return nil, err
	}
	log.Debug("Parsing value %v !\n", valueDecl)

	if quoted {
		log.Debug("consume quote %v\n", p.cur())
		_, err := p.consumeToken(SimpleQuoteToken, DoubleQuoteToken)
		if err != nil {
			debug("uuuh, wasn't a quote")
			return nil, err
		}
	}

	return valueDecl, nil
}

// parseJoin parses the JOIN keywords and all its condition
// JOIN user_addresses ON address.id=user_addresses.address_id
func (p *parser) parseJoin() (*Decl, error) {
	joinDecl, err := p.consumeToken(JoinToken)
	if err != nil {
		return nil, err
	}

	// TABLE NAME
	tableDecl, err := p.parseAttribute()
	if err != nil {
		return nil, err
	}
	joinDecl.Add(tableDecl)

	// ON
	onDecl, err := p.consumeToken(OnToken)
	if err != nil {
		return nil, err
	}
	// onDecl := NewDecl(t)
	joinDecl.Add(onDecl)

	// ATTRIBUTE
	leftAttributeDecl, err := p.parseAttribute()
	if err != nil {
		return nil, err
	}
	onDecl.Add(leftAttributeDecl)

	// EQUAL
	equalAttr, err := p.consumeToken(EqualityToken)
	if err != nil {
		return nil, err
	}
	onDecl.Add(equalAttr)

	// ATTRIBUTE
	rightAttributeDecl, err := p.parseAttribute()
	if err != nil {
		return nil, err
	}
	onDecl.Add(rightAttributeDecl)

	return joinDecl, nil
}

func (p *parser) parseListElement() (*Decl, error) {
	quoted := false

	// In case of INSERT, can be DEFAULT here
	if p.is(DefaultToken) {
		v, err := p.consumeToken(DefaultToken)
		if err != nil {
			return nil, err
		}
		return v, nil
	}

	if p.is(SimpleQuoteToken) || p.is(DoubleQuoteToken) {
		quoted = true
		p.next()
	}

	var valueDecl *Decl
	valueDecl, err := p.consumeToken(StringToken, NumberToken, NullToken, DateToken, NowToken, FalseToken)
	if err != nil {
		return nil, err
	}

	if quoted {
		if _, err := p.consumeToken(SimpleQuoteToken, DoubleQuoteToken); err != nil {
			return nil, err
		}
	}

	return valueDecl, nil
}

func (p *parser) hasNext() bool {
	if p.index+1 < p.tokenLen {
		return true
	}
	return false
}

func (p *parser) next() error {
	if !p.hasNext() {
		return fmt.Errorf("Unexpected end of tokens")
	}

	p.index++
	return nil
}

func (p *parser) peekBackward() Token {
	return p.tokens[p.index-1]
}

func (p *parser) cur() Token {
	return p.tokens[p.index]
}

func (p *parser) peekForward() Token {
	return p.tokens[p.index+1]
}

func (p *parser) is(tokenTypes ...int) bool {

	for _, tokenType := range tokenTypes {
		if p.cur().Token == tokenType {
			return true
		}
	}

	return false
}

func (p *parser) isNot(tokenTypes ...int) bool {
	return !p.is(tokenTypes...)
}

func (p *parser) isNext(tokenTypes ...int) (t Token, err error) {

	if !p.hasNext() {
		debug("parser.isNext: has no next")
		return t, p.syntaxError()
	}

	debug("parser.isNext %v", tokenTypes)
	for _, tokenType := range tokenTypes {
		if p.peekForward().Token == tokenType {
			return p.peekForward(), nil
		}
	}

	debug("parser.isNext: Next (%v) is not among %v", p.cur(), tokenTypes)
	return t, p.syntaxError()
}

func (p *parser) mustHaveNext(tokenTypes ...int) (t Token, err error) {

	if !p.hasNext() {
		debug("parser.mustHaveNext: has no next")
		return t, p.syntaxError()
	}

	if err = p.next(); err != nil {
		debug("parser.mustHaveNext: error getting next")
		return t, err
	}

	debug("parser.mustHaveNext %v", tokenTypes)
	for _, tokenType := range tokenTypes {
		if p.is(tokenType) {
			return p.cur(), nil
		}
	}

	debug("parser.mustHaveNext: Next (%v) is not among %v", p.cur(), tokenTypes)
	return t, p.syntaxError()
}

func (p *parser) consumeToken(tokenTypes ...int) (*Decl, error) {

	if !p.is(tokenTypes...) {
		return nil, p.syntaxError()
	}

	decl := NewDecl(p.cur())
	p.next()
	return decl, nil
}

func (p *parser) syntaxError() error {
	if p.index == 0 {
		return fmt.Errorf("Syntax error near %v %v", p.cur().Lexeme, p.peekForward().Lexeme)
	} else if !p.hasNext() {
		return fmt.Errorf("Syntax error near %v %v", p.peekBackward().Lexeme, p.cur().Lexeme)
	}
	return fmt.Errorf("Syntax error near %v %v %v", p.peekBackward().Lexeme, p.cur().Lexeme, p.peekForward().Lexeme)
}

func stripSpaces(t []Token) []Token {
	retT := []Token{}

	for i := range t {
		if t[i].Token != SpaceToken {
			retT = append(retT, t[i])
		}
	}

	return retT
}
