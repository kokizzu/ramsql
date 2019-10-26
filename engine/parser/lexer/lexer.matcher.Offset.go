package lexer

func (l *Lexer) MatchOffsetToken() bool {
  return l.Match([]byte("offset"), OffsetToken)
}
