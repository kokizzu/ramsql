package parser

func (l *lexer) MatchSetToken() bool {
  return l.Match([]byte("set"), SetToken)
}
