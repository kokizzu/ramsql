package parser

func (l *lexer) MatchGreaterOrEqualToken() bool {
  return l.Match([]byte(">="), GreaterOrEqualToken)
}
