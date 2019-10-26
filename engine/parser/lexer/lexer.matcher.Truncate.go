package lexer

func (l *Lexer) matchTruncateToken() bool {
  return l.match([]byte("truncate"), TruncateToken)
}
