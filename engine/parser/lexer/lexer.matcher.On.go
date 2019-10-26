package lexer

func (l *Lexer) matchOnToken() bool {
  return l.match([]byte("on"), OnToken)
}
