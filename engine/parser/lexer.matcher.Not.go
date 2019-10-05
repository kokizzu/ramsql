package parser

func (l *lexer) MatchNotToken() bool {
  return l.Match([]byte("not"), NotToken)
}
