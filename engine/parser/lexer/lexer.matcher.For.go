package lexer

func (l *Lexer) matchForToken() bool {
  return l.match([]byte("for"), ForToken)
}
