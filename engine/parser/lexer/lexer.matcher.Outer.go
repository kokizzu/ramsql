package lexer

func (l *Lexer) MatchOuterToken() bool {
  return l.Match([]byte("outer"), OuterToken)
}
