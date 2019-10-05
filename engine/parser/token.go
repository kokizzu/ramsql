package parser

// Token contains a token id and its lexeme
// TODO: Add Location{Path string, Line int, Offset int}
type Token struct {
	Token  int
	Lexeme string
}
