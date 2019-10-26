package lexer

func (l *Lexer) MatchFromToken() bool {
  return l.Match([]byte("from"), FromToken)
}
