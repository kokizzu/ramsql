package parser

func (l *lexer) MatchReturningToken() bool {
  return l.Match([]byte("returning"), ReturningToken)
}
