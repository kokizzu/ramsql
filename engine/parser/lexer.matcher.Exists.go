package parser

func (l *lexer) MatchExistsToken() bool {
  return l.Match([]byte("exists"), ExistsToken)
}
