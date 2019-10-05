package parser

func (l *lexer) MatchOnToken() bool {
  return l.Match([]byte("on"), OnToken)
}
