package parser

func (l *lexer) MatchSelectToken() bool {
  return l.Match([]byte("select"), SelectToken)
}
