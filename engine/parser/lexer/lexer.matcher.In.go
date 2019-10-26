package lexer

func (l *Lexer) matchInToken() bool {
  return l.match([]byte("in"), InToken)
}
