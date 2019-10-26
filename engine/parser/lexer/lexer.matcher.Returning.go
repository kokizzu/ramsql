package lexer

func (l *Lexer) MatchReturningToken() bool {
  return l.Match([]byte("returning"), ReturningToken)
}
