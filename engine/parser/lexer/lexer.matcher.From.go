package lexer

func (l *Lexer) matchFromToken() bool {
  return l.match([]byte("from"), FromToken)
}
