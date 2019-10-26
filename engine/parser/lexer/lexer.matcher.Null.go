package lexer

func (l *Lexer) matchNullToken() bool {
  return l.match([]byte("null"), NullToken)
}
