package lexer

func (l *Lexer) matchAscToken() bool {
  return l.match([]byte("asc"), AscToken)
}
