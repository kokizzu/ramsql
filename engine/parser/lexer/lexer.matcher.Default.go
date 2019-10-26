package lexer

func (l *Lexer) matchDefaultToken() bool {
  return l.match([]byte("default"), DefaultToken)
}
