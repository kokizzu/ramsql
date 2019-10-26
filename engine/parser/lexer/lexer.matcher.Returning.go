package lexer

func (l *Lexer) matchReturningToken() bool {
  return l.match([]byte("returning"), ReturningToken)
}
