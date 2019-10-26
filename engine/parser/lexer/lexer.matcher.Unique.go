package lexer

func (l *Lexer) MatchUniqueToken() bool {
  return l.Match([]byte("unique"), UniqueToken)
}
