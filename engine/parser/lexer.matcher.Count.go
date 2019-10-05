package parser

func (l *lexer) MatchCountToken() bool {
  return l.Match([]byte("count"), CountToken)
}
