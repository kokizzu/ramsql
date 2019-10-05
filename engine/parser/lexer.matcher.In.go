package parser

func (l *lexer) MatchInToken() bool {
  return l.Match([]byte("in"), InToken)
}
