package parser

func (l *lexer) MatchDescToken() bool {
  return l.Match([]byte("desc"), DescToken)
}
