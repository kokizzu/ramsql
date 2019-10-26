package lexer

func (l *Lexer) matchCountToken() bool {
  return l.match([]byte("count"), CountToken)
}
