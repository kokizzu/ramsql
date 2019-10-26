package lexer

func (l *Lexer) MatchTruncateToken() bool {
  return l.Match([]byte("truncate"), TruncateToken)
}
