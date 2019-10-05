package parser

func (l *lexer) MatchFromToken() bool {
  return l.Match([]byte("from"), FromToken)
}
