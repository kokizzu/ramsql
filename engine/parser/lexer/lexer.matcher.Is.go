package lexer

func (l *Lexer) matchIsToken() bool {
  return l.match([]byte("is"), IsToken)
}
