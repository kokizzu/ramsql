package parser

func (l *lexer) MatchUsingToken() bool {
  return l.Match([]byte("using"), UsingToken)
}
