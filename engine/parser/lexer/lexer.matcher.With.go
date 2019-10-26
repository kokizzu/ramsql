package lexer

func (l *Lexer) matchWithToken() bool {
  return l.match([]byte("with"), WithToken)
}
