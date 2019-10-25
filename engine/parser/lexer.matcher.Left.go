package parser

func (l *lexer) MatchLeftToken() bool {
  return l.Match([]byte("left"), LeftToken)
}
