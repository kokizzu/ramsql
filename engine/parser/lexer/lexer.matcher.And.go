package lexer

func (l *Lexer) matchAndToken() bool {
  return l.match([]byte("and"), AndToken)
}
