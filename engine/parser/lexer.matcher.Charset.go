package parser

func (l *lexer) MatchCharsetToken() bool {
  return l.Match([]byte("charset"), CharsetToken)
}
