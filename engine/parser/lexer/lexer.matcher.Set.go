package lexer

func (l *Lexer) matchSetToken() bool {
  return l.match([]byte("set"), SetToken)
}
