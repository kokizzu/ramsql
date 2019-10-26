package lexer

func (l *Lexer) matchOuterToken() bool {
  return l.match([]byte("outer"), OuterToken)
}
