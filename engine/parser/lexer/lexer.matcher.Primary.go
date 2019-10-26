package lexer

func (l *Lexer) matchPrimaryToken() bool {
  return l.match([]byte("primary"), PrimaryToken)
}
