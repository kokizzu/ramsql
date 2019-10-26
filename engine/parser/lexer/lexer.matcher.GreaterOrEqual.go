package lexer

func (l *Lexer) MatchGreaterOrEqualToken() bool {
  return l.Match([]byte(">="), GreaterOrEqualToken)
}
