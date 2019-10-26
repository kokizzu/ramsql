package lexer

func (l *Lexer) matchGreaterOrEqualToken() bool {
  return l.match([]byte(">="), GreaterOrEqualToken)
}
