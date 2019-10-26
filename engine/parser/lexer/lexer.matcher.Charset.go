package lexer

func (l *Lexer) MatchCharsetToken() bool {
  return l.Match([]byte("charset"), CharsetToken)
}
