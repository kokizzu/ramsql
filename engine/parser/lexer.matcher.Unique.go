package parser

func (l *lexer) MatchUniqueToken() bool {
  return l.Match([]byte("unique"), UniqueToken)
}
