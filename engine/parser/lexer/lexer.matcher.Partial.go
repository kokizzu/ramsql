package lexer

func (l *Lexer) matchPartialToken() bool {
  return l.match([]byte("partial"), PartialToken)
}
