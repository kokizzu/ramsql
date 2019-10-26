package lexer

func (l *Lexer) matchSelectToken() bool {
  return l.match([]byte("select"), SelectToken)
}
