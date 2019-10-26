package lexer

import (
	"fmt"
	"time"
	"unicode"

	"github.com/mlhoyt/ramsql/engine/log"
)

// Lexer contains the state of the instruction being processed
type Lexer struct {
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

// Lex accepts an instruction and returns the extracted token slice
func (l *Lexer) Lex(instruction []byte) ([]Token, error) {
	l.instructionLen = len(instruction)
	l.tokens = nil
	l.instruction = instruction
	l.pos = 0
	securityPos := 0

	var matchers []Matcher
	// Punctuation Matcher
	matchers = append(matchers, l.matchSpaceToken)
	matchers = append(matchers, l.matchSemicolonToken)
	matchers = append(matchers, l.matchCommaToken)
	matchers = append(matchers, l.matchBracketOpeningToken)
	matchers = append(matchers, l.matchBracketClosingToken)
	matchers = append(matchers, l.matchStarToken)
	matchers = append(matchers, l.matchSimpleQuoteToken)
	matchers = append(matchers, l.matchEqualityToken)
	matchers = append(matchers, l.matchPeriodToken)
	matchers = append(matchers, l.matchDoubleQuoteToken)
	matchers = append(matchers, l.matchLessOrEqualToken)
	matchers = append(matchers, l.matchLeftDipleToken)
	matchers = append(matchers, l.matchGreaterOrEqualToken)
	matchers = append(matchers, l.matchRightDipleToken)
	matchers = append(matchers, l.matchBacktickToken)
	// First order Matcher
	matchers = append(matchers, l.matchCreateToken)
	matchers = append(matchers, l.matchDeleteToken)
	matchers = append(matchers, l.matchDropToken)
	matchers = append(matchers, l.matchGrantToken)
	matchers = append(matchers, l.matchInsertToken)
	matchers = append(matchers, l.matchSelectToken)
	matchers = append(matchers, l.matchTruncateToken)
	matchers = append(matchers, l.matchUpdateToken)
	// Second order Matcher
	matchers = append(matchers, l.matchActionToken)
	matchers = append(matchers, l.matchAndToken)
	matchers = append(matchers, l.matchAscToken)
	matchers = append(matchers, l.matchAsToken)
	matchers = append(matchers, l.matchAutoincrementToken)
	matchers = append(matchers, l.matchBtreeToken)
	matchers = append(matchers, l.matchByToken)
	matchers = append(matchers, l.matchCascadeToken)
	matchers = append(matchers, l.matchCharacterToken)
	matchers = append(matchers, l.matchCharsetToken)
	matchers = append(matchers, l.matchConstraintToken)
	matchers = append(matchers, l.matchCountToken)
	matchers = append(matchers, l.matchDefaultToken)
	matchers = append(matchers, l.matchDescToken)
	matchers = append(matchers, l.matchEngineToken)
	matchers = append(matchers, l.matchExistsToken)
	matchers = append(matchers, l.matchFalseToken)
	matchers = append(matchers, l.matchForeignToken)
	matchers = append(matchers, l.matchForToken)
	matchers = append(matchers, l.matchFromToken)
	matchers = append(matchers, l.matchFullToken)
	matchers = append(matchers, l.matchHashToken)
	matchers = append(matchers, l.matchIfToken)
	matchers = append(matchers, l.matchIndexToken)
	matchers = append(matchers, l.matchInnerToken)
	matchers = append(matchers, l.matchIntoToken)
	matchers = append(matchers, l.matchInToken)
	matchers = append(matchers, l.matchIsToken)
	matchers = append(matchers, l.matchJoinToken)
	matchers = append(matchers, l.matchKeyToken)
	matchers = append(matchers, l.matchLeftToken)
	matchers = append(matchers, l.matchLimitToken)
	matchers = append(matchers, l.matchLocalTimestampToken)
	matchers = append(matchers, l.matchMatchToken)
	matchers = append(matchers, l.matchNotToken)
	matchers = append(matchers, l.matchNowToken)
	matchers = append(matchers, l.matchNoToken)
	matchers = append(matchers, l.matchNullToken)
	matchers = append(matchers, l.matchOffsetToken)
	matchers = append(matchers, l.matchOnToken)
	matchers = append(matchers, l.matchOrderToken)
	matchers = append(matchers, l.matchOrToken)
	matchers = append(matchers, l.matchOuterToken)
	matchers = append(matchers, l.matchPartialToken)
	matchers = append(matchers, l.matchPrimaryToken)
	matchers = append(matchers, l.matchReferencesToken)
	matchers = append(matchers, l.matchRestrictToken)
	matchers = append(matchers, l.matchReturningToken)
	matchers = append(matchers, l.matchRightToken)
	matchers = append(matchers, l.matchSetToken)
	matchers = append(matchers, l.matchSimpleToken)
	matchers = append(matchers, l.matchTableToken)
	matchers = append(matchers, l.matchTimeToken)
	matchers = append(matchers, l.matchUniqueToken)
	matchers = append(matchers, l.matchUsingToken)
	matchers = append(matchers, l.matchValuesToken)
	matchers = append(matchers, l.matchWhereToken)
	matchers = append(matchers, l.matchWithToken)
	matchers = append(matchers, l.matchZoneToken)
	// Type Matcher
	matchers = append(matchers, l.matchEscapedStringToken)
	matchers = append(matchers, l.matchDateToken)
	matchers = append(matchers, l.matchNumberToken)
	matchers = append(matchers, l.matchStringToken)

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

func (l *Lexer) matchSpaceToken() bool {

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

func (l *Lexer) matchStringToken() bool {

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

func (l *Lexer) matchNumberToken() bool {

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

// matchDateToken prefers time.RFC3339Nano but will match a few others as well
func (l *Lexer) matchDateToken() bool {

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

func (l *Lexer) matchDoubleQuoteToken() bool {

	if l.instruction[l.pos] == '"' {

		t := Token{
			Token:  DoubleQuoteToken,
			Lexeme: "\"",
		}
		l.tokens = append(l.tokens, t)
		l.pos++

		if l.matchDoubleQuotedStringToken() {
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

func (l *Lexer) matchEscapedStringToken() bool {
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

func (l *Lexer) matchDoubleQuotedStringToken() bool {
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

func (l *Lexer) matchSimpleQuoteToken() bool {

	if l.instruction[l.pos] == '\'' {

		t := Token{
			Token:  SimpleQuoteToken,
			Lexeme: "'",
		}
		l.tokens = append(l.tokens, t)
		l.pos++

		if l.matchSingleQuotedStringToken() {
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

func (l *Lexer) matchSingleQuotedStringToken() bool {
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

func (l *Lexer) matchSingle(char byte, token int) bool {

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

func (l *Lexer) match(str []byte, token int) bool {

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

// DateLongFormat is same as time.RFC3339Nano
const DateLongFormat = time.RFC3339Nano

// DateShortFormat is a short date format with human-readable month element
const DateShortFormat = "2006-Jan-02"

// DateNumberFormat is a fully numeric short date format
const DateNumberFormat = "2006-01-02"

// ParseDate intends to parse all SQL date formats
func ParseDate(data string) (*time.Time, error) {
	t, err := time.Parse(DateLongFormat, data)
	if err == nil {
		return &t, nil
	}

	t, err = time.Parse(time.RFC3339, data)
	if err == nil {
		return &t, nil
	}

	t, err = time.Parse("2006-01-02 15:04:05.999999 -0700 MST", data)
	if err == nil {
		return &t, nil
	}

	t, err = time.Parse(DateShortFormat, data)
	if err == nil {
		return &t, nil
	}

	t, err = time.Parse(DateNumberFormat, data)
	if err == nil {
		return &t, nil
	}

	return nil, fmt.Errorf("not a date")
}
