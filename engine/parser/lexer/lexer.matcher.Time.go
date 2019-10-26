package lexer

func (l *Lexer) MatchTimeToken() bool {
  return l.Match([]byte("time"), TimeToken)
}
