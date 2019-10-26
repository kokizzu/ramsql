package lexer

func (l *Lexer) matchOrToken() bool {
  return l.match([]byte("or"), OrToken)
}
