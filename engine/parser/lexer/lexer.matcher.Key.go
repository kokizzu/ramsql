package lexer

func (l *Lexer) matchKeyToken() bool {
  return l.match([]byte("key"), KeyToken)
}
