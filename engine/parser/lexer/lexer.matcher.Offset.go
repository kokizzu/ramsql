package lexer

func (l *Lexer) matchOffsetToken() bool {
  return l.match([]byte("offset"), OffsetToken)
}
