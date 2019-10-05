package parser

func (l *lexer) MatchOffsetToken() bool {
  return l.Match([]byte("offset"), OffsetToken)
}
