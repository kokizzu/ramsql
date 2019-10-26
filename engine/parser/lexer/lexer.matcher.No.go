package lexer

func (l *Lexer) matchNoToken() bool {
  return l.match([]byte("no"), NoToken)
}
