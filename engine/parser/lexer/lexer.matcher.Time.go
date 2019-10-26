package lexer

func (l *Lexer) matchTimeToken() bool {
  return l.match([]byte("time"), TimeToken)
}
