package lexer

func (l *Lexer) MatchExistsToken() bool {
  return l.Match([]byte("exists"), ExistsToken)
}
