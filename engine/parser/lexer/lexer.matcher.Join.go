package lexer

func (l *Lexer) matchJoinToken() bool {
  return l.match([]byte("join"), JoinToken)
}
