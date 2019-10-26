package lexer

func (l *Lexer) matchUsingToken() bool {
  return l.match([]byte("using"), UsingToken)
}
