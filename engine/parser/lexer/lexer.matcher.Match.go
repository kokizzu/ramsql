package lexer

func (l *Lexer) MatchMatchToken() bool {
  return l.Match([]byte("match"), MatchToken)
}
