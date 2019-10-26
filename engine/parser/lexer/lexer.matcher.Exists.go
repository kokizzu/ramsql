package lexer

func (l *Lexer) matchExistsToken() bool {
  return l.match([]byte("exists"), ExistsToken)
}
