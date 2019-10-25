package parser

import (
	"fmt"
	"unicode"

	"github.com/mlhoyt/ramsql/engine/log"
)

type lexer struct {
	tokens         []Token
	instruction    []byte
	instructionLen int
	pos            int
}

// SQL Tokens
const (
	ActionToken         = iota // Second-order
	AndToken                   // Second-order
	AsToken                    // Second-order
	AscToken                   // Second-order
	AutoincrementToken         // Second-order
	BacktickToken              // Punctuation
	BracketClosingToken        // Punctuation
	BracketOpeningToken        // Punctuation
	BtreeToken                 // Second-order
	ByToken                    // Second-order
	CascadeToken               // Second-order
	CharacterToken             // Second-order
	CharsetToken               // Second-order
	CommaToken                 // Punctuation
	ConstraintToken            // Second-order
	CountToken                 // Second-order
	CreateToken                // First-order
	DateToken                  // Type
	DefaultToken               // Second-order
	DeleteToken                // First-order
	DescToken                  // Second-order
	DoubleQuoteToken           // Quote
	DropToken                  // First-order
	EngineToken                // Second-order
	EqualityToken              // Quote
	ExistsToken                // Second-order
	ExplainToken               // First-order
	FalseToken                 // Second-order
	ForToken                   // Second-order
	ForeignToken               // Second-order
	FromToken                  // Second-order
	FullToken                  // Second-order
	GrantToken                 // First-order
	GreaterOrEqualToken        // Punctuation
	HashToken                  // Second-order
	IfToken                    // Second-order
	InToken                    // Second-order
	IndexToken                 // Second-order
	InnerToken                 // Second-order
	InsertToken                // First-order
	IntToken                   // Type
	IntoToken                  // Second-order
	IsToken                    // Second-order
	JoinToken                  // Second-order
	KeyToken                   // Type
	LeftToken                  // Second-order
	LeftDipleToken             // Punctuation
	LessOrEqualToken           // Punctuation
	LimitToken                 // Second-order
	LocalTimestampToken        // Second-order
	MatchToken                 // Second-order
	NoToken                    // Second-order
	NotToken                   // Second-order
	NowToken                   // Second-order
	NullToken                  // Second-order
	NumberToken                // Type
	OffsetToken                // Second-order
	OnToken                    // Second-order
	OrToken                    // Second-order
	OrderToken                 // Second-order
	OuterToken                 // Second-order
	PartialToken               // Quote
	PeriodToken                // Quote
	PrimaryToken               // Type
	ReferencesToken            // Second-order
	ReturningToken             // Second-order
	RestrictToken              // Second-order
	RightToken                 // Second-order
	RightDipleToken            // Punctuation
	SelectToken                // First-order
	SemicolonToken             // Punctuation
	SetToken                   // Second-order
	SimpleToken                // Second-order
	SimpleQuoteToken           // Quote
	SpaceToken                 // Punctuation
	StarToken                  // Quote
	StringToken                // Type
	TableToken                 // Second-order
	TextToken                  // Type
	TimeToken                  // Second-order
	TruncateToken              // First-order
	UniqueToken                // Second-order
	UpdateToken                // First-order
	UsingToken                 // Second-order
	ValuesToken                // Second-order
	WhereToken                 // Second-order
	WithToken                  // Second-order
	ZoneToken                  // Second-order
)

// Matcher tries to match given string to an SQL token
type Matcher func() bool

//go:generate ./lexer-generate-matcher.sh --lexeme "(" --name BracketOpening
//go:generate ./lexer-generate-matcher.sh --lexeme ")" --name BracketClosing
//go:generate ./lexer-generate-matcher.sh --lexeme "*" --name Star
//go:generate ./lexer-generate-matcher.sh --lexeme "," --name Comma
//go:generate ./lexer-generate-matcher.sh --lexeme "." --name Period
//go:generate ./lexer-generate-matcher.sh --lexeme ";" --name Semicolon
//go:generate ./lexer-generate-matcher.sh --lexeme "<" --name LeftDiple
//go:generate ./lexer-generate-matcher.sh --lexeme "<=" --name LessOrEqual
//go:generate ./lexer-generate-matcher.sh --lexeme "=" --name Equality
//go:generate ./lexer-generate-matcher.sh --lexeme ">" --name RightDiple
//go:generate ./lexer-generate-matcher.sh --lexeme ">=" --name GreaterOrEqual
//go:generate ./lexer-generate-matcher.sh --lexeme "`" --name Backtick
//go:generate ./lexer-generate-matcher.sh --lexeme "action"
//go:generate ./lexer-generate-matcher.sh --lexeme "and"
//go:generate ./lexer-generate-matcher.sh --lexeme "as"
//go:generate ./lexer-generate-matcher.sh --lexeme "asc"
//go:generate ./lexer-generate-matcher.sh --lexeme "autoincrement" --lexeme "auto_increment"
//go:generate ./lexer-generate-matcher.sh --lexeme "btree"
//go:generate ./lexer-generate-matcher.sh --lexeme "by"
//go:generate ./lexer-generate-matcher.sh --lexeme "cascade"
//go:generate ./lexer-generate-matcher.sh --lexeme "character"
//go:generate ./lexer-generate-matcher.sh --lexeme "charset"
//go:generate ./lexer-generate-matcher.sh --lexeme "constraint"
//go:generate ./lexer-generate-matcher.sh --lexeme "count"
//go:generate ./lexer-generate-matcher.sh --lexeme "create"
//go:generate ./lexer-generate-matcher.sh --lexeme "default"
//go:generate ./lexer-generate-matcher.sh --lexeme "delete"
//go:generate ./lexer-generate-matcher.sh --lexeme "desc"
//go:generate ./lexer-generate-matcher.sh --lexeme "drop"
//go:generate ./lexer-generate-matcher.sh --lexeme "engine"
//go:generate ./lexer-generate-matcher.sh --lexeme "exists"
//go:generate ./lexer-generate-matcher.sh --lexeme "false"
//go:generate ./lexer-generate-matcher.sh --lexeme "for"
//go:generate ./lexer-generate-matcher.sh --lexeme "foreign"
//go:generate ./lexer-generate-matcher.sh --lexeme "from"
//go:generate ./lexer-generate-matcher.sh --lexeme "full"
//go:generate ./lexer-generate-matcher.sh --lexeme "grant"
//go:generate ./lexer-generate-matcher.sh --lexeme "hash"
//go:generate ./lexer-generate-matcher.sh --lexeme "if"
//go:generate ./lexer-generate-matcher.sh --lexeme "in"
//go:generate ./lexer-generate-matcher.sh --lexeme "index"
//go:generate ./lexer-generate-matcher.sh --lexeme "inner"
//go:generate ./lexer-generate-matcher.sh --lexeme "insert"
//go:generate ./lexer-generate-matcher.sh --lexeme "into"
//go:generate ./lexer-generate-matcher.sh --lexeme "is"
//go:generate ./lexer-generate-matcher.sh --lexeme "join"
//go:generate ./lexer-generate-matcher.sh --lexeme "key"
//go:generate ./lexer-generate-matcher.sh --lexeme "left"
//go:generate ./lexer-generate-matcher.sh --lexeme "limit"
//go:generate ./lexer-generate-matcher.sh --lexeme "localtimestamp" --lexeme "current_timestamp" --name LocalTimestamp
//go:generate ./lexer-generate-matcher.sh --lexeme "match"
//go:generate ./lexer-generate-matcher.sh --lexeme "no"
//go:generate ./lexer-generate-matcher.sh --lexeme "not"
//go:generate ./lexer-generate-matcher.sh --lexeme "now()" --name Now
//go:generate ./lexer-generate-matcher.sh --lexeme "null"
//go:generate ./lexer-generate-matcher.sh --lexeme "offset"
//go:generate ./lexer-generate-matcher.sh --lexeme "on"
//go:generate ./lexer-generate-matcher.sh --lexeme "or"
//go:generate ./lexer-generate-matcher.sh --lexeme "order"
//go:generate ./lexer-generate-matcher.sh --lexeme "outer"
//go:generate ./lexer-generate-matcher.sh --lexeme "partial"
//go:generate ./lexer-generate-matcher.sh --lexeme "primary"
//go:generate ./lexer-generate-matcher.sh --lexeme "references"
//go:generate ./lexer-generate-matcher.sh --lexeme "restrict"
//go:generate ./lexer-generate-matcher.sh --lexeme "returning"
//go:generate ./lexer-generate-matcher.sh --lexeme "right"
//go:generate ./lexer-generate-matcher.sh --lexeme "select"
//go:generate ./lexer-generate-matcher.sh --lexeme "set"
//go:generate ./lexer-generate-matcher.sh --lexeme "simple"
//go:generate ./lexer-generate-matcher.sh --lexeme "table"
//go:generate ./lexer-generate-matcher.sh --lexeme "time"
//go:generate ./lexer-generate-matcher.sh --lexeme "truncate"
//go:generate ./lexer-generate-matcher.sh --lexeme "unique"
//go:generate ./lexer-generate-matcher.sh --lexeme "update"
//go:generate ./lexer-generate-matcher.sh --lexeme "using"
//go:generate ./lexer-generate-matcher.sh --lexeme "values"
//go:generate ./lexer-generate-matcher.sh --lexeme "where"
//go:generate ./lexer-generate-matcher.sh --lexeme "with"
//go:generate ./lexer-generate-matcher.sh --lexeme "zone"

func (l *lexer) lex(instruction []byte) ([]Token, error) {
	l.instructionLen = len(instruction)
	l.tokens = nil
	l.instruction = instruction
	l.pos = 0
	securityPos := 0

	var matchers []Matcher
	// Punctuation Matcher
	matchers = append(matchers, l.MatchSpaceToken)
	matchers = append(matchers, l.MatchSemicolonToken)
	matchers = append(matchers, l.MatchCommaToken)
	matchers = append(matchers, l.MatchBracketOpeningToken)
	matchers = append(matchers, l.MatchBracketClosingToken)
	matchers = append(matchers, l.MatchStarToken)
	matchers = append(matchers, l.MatchSimpleQuoteToken)
	matchers = append(matchers, l.MatchEqualityToken)
	matchers = append(matchers, l.MatchPeriodToken)
	matchers = append(matchers, l.MatchDoubleQuoteToken)
	matchers = append(matchers, l.MatchLessOrEqualToken)
	matchers = append(matchers, l.MatchLeftDipleToken)
	matchers = append(matchers, l.MatchGreaterOrEqualToken)
	matchers = append(matchers, l.MatchRightDipleToken)
	matchers = append(matchers, l.MatchBacktickToken)
	// First order Matcher
	matchers = append(matchers, l.MatchCreateToken)
	matchers = append(matchers, l.MatchDeleteToken)
	matchers = append(matchers, l.MatchDropToken)
	matchers = append(matchers, l.MatchGrantToken)
	matchers = append(matchers, l.MatchInsertToken)
	matchers = append(matchers, l.MatchSelectToken)
	matchers = append(matchers, l.MatchTruncateToken)
	matchers = append(matchers, l.MatchUpdateToken)
	// Second order Matcher
	matchers = append(matchers, l.MatchActionToken)
	matchers = append(matchers, l.MatchAndToken)
	matchers = append(matchers, l.MatchAscToken)
	matchers = append(matchers, l.MatchAsToken)
	matchers = append(matchers, l.MatchAutoincrementToken)
	matchers = append(matchers, l.MatchBtreeToken)
	matchers = append(matchers, l.MatchByToken)
	matchers = append(matchers, l.MatchCascadeToken)
	matchers = append(matchers, l.MatchCharacterToken)
	matchers = append(matchers, l.MatchCharsetToken)
	matchers = append(matchers, l.MatchConstraintToken)
	matchers = append(matchers, l.MatchCountToken)
	matchers = append(matchers, l.MatchDefaultToken)
	matchers = append(matchers, l.MatchDescToken)
	matchers = append(matchers, l.MatchEngineToken)
	matchers = append(matchers, l.MatchExistsToken)
	matchers = append(matchers, l.MatchFalseToken)
	matchers = append(matchers, l.MatchForeignToken)
	matchers = append(matchers, l.MatchForToken)
	matchers = append(matchers, l.MatchFromToken)
	matchers = append(matchers, l.MatchFullToken)
	matchers = append(matchers, l.MatchHashToken)
	matchers = append(matchers, l.MatchIfToken)
	matchers = append(matchers, l.MatchIndexToken)
	matchers = append(matchers, l.MatchInnerToken)
	matchers = append(matchers, l.MatchIntoToken)
	matchers = append(matchers, l.MatchInToken)
	matchers = append(matchers, l.MatchIsToken)
	matchers = append(matchers, l.MatchJoinToken)
	matchers = append(matchers, l.MatchKeyToken)
	matchers = append(matchers, l.MatchLeftToken)
	matchers = append(matchers, l.MatchLimitToken)
	matchers = append(matchers, l.MatchLocalTimestampToken)
	matchers = append(matchers, l.MatchMatchToken)
	matchers = append(matchers, l.MatchNotToken)
	matchers = append(matchers, l.MatchNowToken)
	matchers = append(matchers, l.MatchNoToken)
	matchers = append(matchers, l.MatchNullToken)
	matchers = append(matchers, l.MatchOffsetToken)
	matchers = append(matchers, l.MatchOnToken)
	matchers = append(matchers, l.MatchOrderToken)
	matchers = append(matchers, l.MatchOrToken)
	matchers = append(matchers, l.MatchOuterToken)
	matchers = append(matchers, l.MatchPartialToken)
	matchers = append(matchers, l.MatchPrimaryToken)
	matchers = append(matchers, l.MatchReferencesToken)
	matchers = append(matchers, l.MatchRestrictToken)
	matchers = append(matchers, l.MatchReturningToken)
	matchers = append(matchers, l.MatchRightToken)
	matchers = append(matchers, l.MatchSetToken)
	matchers = append(matchers, l.MatchSimpleToken)
	matchers = append(matchers, l.MatchTableToken)
	matchers = append(matchers, l.MatchTimeToken)
	matchers = append(matchers, l.MatchUniqueToken)
	matchers = append(matchers, l.MatchUsingToken)
	matchers = append(matchers, l.MatchValuesToken)
	matchers = append(matchers, l.MatchWhereToken)
	matchers = append(matchers, l.MatchWithToken)
	matchers = append(matchers, l.MatchZoneToken)
	// Type Matcher
	matchers = append(matchers, l.MatchEscapedStringToken)
	matchers = append(matchers, l.MatchDateToken)
	matchers = append(matchers, l.MatchNumberToken)
	matchers = append(matchers, l.MatchStringToken)

	var r bool
	for l.pos < l.instructionLen {
		// fmt.Printf("Tokens : %v\n\n", l.tokens)

		r = false
		for _, m := range matchers {
			if r = m(); r == true {
				securityPos = l.pos
				break
			}
		}

		if r {
			continue
		}

		if l.pos == securityPos {
			log.Warning("Cannot lex <%s>, stuck at pos %d -> [%c]", l.instruction, l.pos, l.instruction[l.pos])
			return nil, fmt.Errorf("Cannot lex instruction. Syntax error near %s", instruction[l.pos:])
		}
		securityPos = l.pos
	}

	return l.tokens, nil
}

func (l *lexer) MatchSpaceToken() bool {

	if unicode.IsSpace(rune(l.instruction[l.pos])) {
		t := Token{
			Token:  SpaceToken,
			Lexeme: " ",
		}
		l.tokens = append(l.tokens, t)
		l.pos++
		return true
	}

	return false
}

func (l *lexer) MatchStringToken() bool {

	i := l.pos
	for i < l.instructionLen &&
		(unicode.IsLetter(rune(l.instruction[i])) ||
			unicode.IsDigit(rune(l.instruction[i])) ||
			l.instruction[i] == '_' ||
			l.instruction[i] == '@' /* || l.instruction[i] == '.'*/) {
		i++
	}

	if i != l.pos {
		t := Token{
			Token:  StringToken,
			Lexeme: string(l.instruction[l.pos:i]),
		}
		l.tokens = append(l.tokens, t)
		l.pos = i
		return true
	}

	return false
}

func (l *lexer) MatchNumberToken() bool {

	i := l.pos
	for i < l.instructionLen && unicode.IsDigit(rune(l.instruction[i])) {
		i++
	}

	if i != l.pos {
		t := Token{
			Token:  NumberToken,
			Lexeme: string(l.instruction[l.pos:i]),
		}
		l.tokens = append(l.tokens, t)
		l.pos = i
		return true
	}

	return false
}

// MatchDateToken prefers time.RFC3339Nano but will match a few others as well
func (l *lexer) MatchDateToken() bool {

	i := l.pos
	for i < l.instructionLen &&
		l.instruction[i] != ',' &&
		l.instruction[i] != ')' {
		i++
	}

	data := string(l.instruction[l.pos:i])

	_, err := ParseDate(data)
	if err != nil {
		return false
	}

	t := Token{
		Token:  StringToken,
		Lexeme: data,
	}

	l.tokens = append(l.tokens, t)
	l.pos = i
	return true
}

func (l *lexer) MatchDoubleQuoteToken() bool {

	if l.instruction[l.pos] == '"' {

		t := Token{
			Token:  DoubleQuoteToken,
			Lexeme: "\"",
		}
		l.tokens = append(l.tokens, t)
		l.pos++

		if l.MatchDoubleQuotedStringToken() {
			t := Token{
				Token:  DoubleQuoteToken,
				Lexeme: "\"",
			}
			l.tokens = append(l.tokens, t)
			l.pos++
			return true
		}

		return true
	}

	return false
}

func (l *lexer) MatchEscapedStringToken() bool {
	i := l.pos
	if l.instruction[i] != '$' || l.instruction[i+1] != '$' {
		return false
	}
	i += 2

	for i+1 < l.instructionLen && !(l.instruction[i] == '$' && l.instruction[i+1] == '$') {
		i++
	}
	i++

	if i == l.instructionLen {
		return false
	}

	tok := NumberToken
	escaped := l.instruction[l.pos+2 : i-1]

	for _, r := range escaped {
		if unicode.IsDigit(rune(r)) == false {
			tok = StringToken
		}
	}

	_, err := ParseDate(string(escaped))
	if err == nil {
		tok = DateToken
	}

	t := Token{
		Token:  tok,
		Lexeme: string(escaped),
	}
	l.tokens = append(l.tokens, t)
	l.pos = i + 1

	return true
}

func (l *lexer) MatchDoubleQuotedStringToken() bool {
	i := l.pos
	for i < l.instructionLen && l.instruction[i] != '"' {
		i++
	}

	t := Token{
		Token:  StringToken,
		Lexeme: string(l.instruction[l.pos:i]),
	}
	l.tokens = append(l.tokens, t)
	l.pos = i

	return true
}

func (l *lexer) MatchSimpleQuoteToken() bool {

	if l.instruction[l.pos] == '\'' {

		t := Token{
			Token:  SimpleQuoteToken,
			Lexeme: "'",
		}
		l.tokens = append(l.tokens, t)
		l.pos++

		if l.MatchSingleQuotedStringToken() {
			t := Token{
				Token:  SimpleQuoteToken,
				Lexeme: "'",
			}
			l.tokens = append(l.tokens, t)
			l.pos++
			return true
		}

		return true
	}

	return false
}

func (l *lexer) MatchSingleQuotedStringToken() bool {
	i := l.pos
	for i < l.instructionLen && l.instruction[i] != '\'' {
		i++
	}

	t := Token{
		Token:  StringToken,
		Lexeme: string(l.instruction[l.pos:i]),
	}
	l.tokens = append(l.tokens, t)
	l.pos = i

	return true
}

func (l *lexer) MatchSingle(char byte, token int) bool {

	if l.pos > l.instructionLen {
		return false
	}

	if l.instruction[l.pos] != char {
		return false
	}

	t := Token{
		Token:  token,
		Lexeme: string(char),
	}

	l.tokens = append(l.tokens, t)
	l.pos++
	return true
}

func (l *lexer) Match(str []byte, token int) bool {

	if l.pos+len(str)-1 > l.instructionLen {
		return false
	}

	// Check for lowercase and uppercase
	for i := range str {
		if unicode.ToLower(rune(l.instruction[l.pos+i])) != unicode.ToLower(rune(str[i])) {
			return false
		}
	}

	// if next character is still a string, it means it doesn't match
	// ie: COUNT shoulnd match COUNTRY
	if l.instructionLen > l.pos+len(str) {
		if unicode.IsLetter(rune(l.instruction[l.pos+len(str)])) ||
			l.instruction[l.pos+len(str)] == '_' {
			return false
		}
	}

	t := Token{
		Token:  token,
		Lexeme: string(str),
	}

	l.tokens = append(l.tokens, t)
	l.pos += len(t.Lexeme)
	return true
}
