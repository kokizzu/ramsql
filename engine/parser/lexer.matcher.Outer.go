package parser

func (l *lexer) MatchOuterToken() bool {
  return l.Match([]byte("outer"), OuterToken)
}
