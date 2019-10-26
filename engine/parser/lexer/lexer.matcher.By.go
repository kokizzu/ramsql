package lexer

func (l *Lexer) matchByToken() bool {
  return l.match([]byte("by"), ByToken)
}
