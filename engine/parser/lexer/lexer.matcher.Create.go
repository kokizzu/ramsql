package lexer

func (l *Lexer) matchCreateToken() bool {
  return l.match([]byte("create"), CreateToken)
}
