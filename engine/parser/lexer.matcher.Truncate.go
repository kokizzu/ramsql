package parser

func (l *lexer) MatchTruncateToken() bool {
  return l.Match([]byte("truncate"), TruncateToken)
}
