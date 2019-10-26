package lexer

func (l *Lexer) matchNotToken() bool {
  return l.match([]byte("not"), NotToken)
}
