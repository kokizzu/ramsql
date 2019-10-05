package parser

func (l *lexer) MatchMatchToken() bool {
  return l.Match([]byte("match"), MatchToken)
}
