package parser

func (l *lexer) MatchTimeToken() bool {
  return l.Match([]byte("time"), TimeToken)
}
