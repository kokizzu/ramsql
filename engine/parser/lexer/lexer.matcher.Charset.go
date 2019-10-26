package lexer

func (l *Lexer) matchCharsetToken() bool {
  return l.match([]byte("charset"), CharsetToken)
}
