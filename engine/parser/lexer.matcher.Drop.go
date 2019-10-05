package parser

func (l *lexer) MatchDropToken() bool {
  return l.Match([]byte("drop"), DropToken)
}
