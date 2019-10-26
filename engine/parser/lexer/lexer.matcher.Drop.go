package lexer

func (l *Lexer) matchDropToken() bool {
  return l.match([]byte("drop"), DropToken)
}
