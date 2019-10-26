package lexer

func (l *Lexer) matchMatchToken() bool {
  return l.match([]byte("match"), MatchToken)
}
