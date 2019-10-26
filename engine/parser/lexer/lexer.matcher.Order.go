package lexer

func (l *Lexer) matchOrderToken() bool {
  return l.match([]byte("order"), OrderToken)
}
