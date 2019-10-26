package lexer

func (l *Lexer) matchActionToken() bool {
  return l.match([]byte("action"), ActionToken)
}
