package parser

import (
	"fmt"
	"github.com/mlhoyt/ramsql/engine/parser/lexer"
)

func (p *parser) parseSelect() (*Instruction, error) {
	i := &Instruction{}
	var err error

	// Create select decl
	selectDecl := NewDecl(p.cur())
	i.Decls = append(i.Decls, selectDecl)

	// After select token, should be either
	// a lexer.StarToken
	// a list of table names + (lexer.StarToken Or Attribute)
	// a builtin func (COUNT, MAX, ...)
	if err = p.next(); err != nil {
		return nil, fmt.Errorf("SELECT token must be followed by attributes to select")
	}

	for {
		if p.is(lexer.CountToken) {
			attrDecl, err := p.parseBuiltinFunc()
			if err != nil {
				return nil, err
			}
			selectDecl.Add(attrDecl)
		} else {
			attrDecl, err := p.parseAttribute()
			if err != nil {
				return nil, err
			}
			selectDecl.Add(attrDecl)
		}

		// If comma, loop again.
		if p.is(lexer.CommaToken) {
			if err := p.next(); err != nil {
				return nil, err
			}
			continue
		}
		break
	}

	// Must be from now
	if p.cur().Token != lexer.FromToken {
		return nil, fmt.Errorf("Syntax error near %v", p.cur())
	}
	fromDecl := NewDecl(p.cur())
	selectDecl.Add(fromDecl)

	// Now must be a list of table
	for {
		// string
		if err = p.next(); err != nil {
			return nil, fmt.Errorf("Unexpected end. Syntax error near %v", p.cur())
		}
		tableNameDecl, err := p.parseAttribute()
		if err != nil {
			return nil, err
		}
		fromDecl.Add(tableNameDecl)

		// If no next, then it's implicit where
		if !p.hasNext() {
			addImplicitWhereAll(selectDecl)
			return i, nil
		}
		// if not comma, break
		if p.cur().Token != lexer.CommaToken {
			break // No more table
		}
	}

	// Optional: INNER
	// var innerJoinDecl *Decl
	if p.is(lexer.InnerToken) {
		_, err := p.consumeToken(lexer.InnerToken)
		if err != nil {
			return nil, err
		}
	}

	// Optional: LEFT, RIGHT
	// var dirOuterJoinDecl *Decl
	if p.is(lexer.LeftToken) {
		_, err := p.consumeToken(lexer.LeftToken)
		if err != nil {
			return nil, err
		}
	} else if p.is(lexer.RightToken) {
		_, err := p.consumeToken(lexer.RightToken)
		if err != nil {
			return nil, err
		}
	}

	// Optional: OUTER
	// var outerJoinDecl *Decl
	if p.is(lexer.OuterToken) {
		_, err := p.consumeToken(lexer.OuterToken)
		if err != nil {
			return nil, err
		}
	}

	// JOIN OR ...?
	for p.is(lexer.JoinToken) {
		joinDecl, err := p.parseJoin()
		if err != nil {
			return nil, err
		}
		// FIXME: Need to annotate joinDecl with innerJoinDecl or outerJoinDecl (and dirOuterJoinDecl)
		selectDecl.Add(joinDecl)
	}

	// Optional: WHERE ..., ORDER [BY] ..., LIMIT ..., OFFSET ..., FOR ...
	hazWhereClause := false
	for {
		switch p.cur().Token {
		case lexer.WhereToken:
			err := p.parseWhere(selectDecl)
			if err != nil {
				return nil, err
			}
			hazWhereClause = true
		case lexer.OrderToken:
			if hazWhereClause == false {
				// WHERE clause is implicit
				addImplicitWhereAll(selectDecl)
			}
			err := p.parseOrderBy(selectDecl)
			if err != nil {
				return nil, err
			}
		case lexer.LimitToken:
			limitDecl, err := p.consumeToken(lexer.LimitToken)
			if err != nil {
				return nil, err
			}
			selectDecl.Add(limitDecl)
			numDecl, err := p.consumeToken(lexer.NumberToken)
			if err != nil {
				return nil, err
			}
			limitDecl.Add(numDecl)
		case lexer.OffsetToken:
			offsetDecl, err := p.consumeToken(lexer.OffsetToken)
			if err != nil {
				return nil, err
			}
			selectDecl.Add(offsetDecl)
			offsetValue, err := p.consumeToken(lexer.NumberToken)
			if err != nil {
				return nil, err
			}
			offsetDecl.Add(offsetValue)
		case lexer.ForToken:
			err := p.parseForUpdate(selectDecl)
			if err != nil {
				return nil, err
			}
		default:
			return i, nil
		}
	}
}

func addImplicitWhereAll(decl *Decl) {

	whereDecl := &Decl{
		Token:  lexer.WhereToken,
		Lexeme: "where",
	}
	whereDecl.Add(&Decl{
		Token:  lexer.NumberToken,
		Lexeme: "1",
	})

	decl.Add(whereDecl)
}

func (p *parser) parseForUpdate(decl *Decl) error {
	// Optionnal
	if !p.is(lexer.ForToken) {
		return nil
	}

	d, err := p.consumeToken(lexer.ForToken)
	if err != nil {
		return err
	}

	u, err := p.consumeToken(lexer.UpdateToken)
	if err != nil {
		return err
	}

	d.Add(u)
	decl.Add(d)
	return nil
}
