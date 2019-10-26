package lexer

func (l *Lexer) matchSimpleToken() bool {
  return l.match([]byte("simple"), SimpleToken)
}
