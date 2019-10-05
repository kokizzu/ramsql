package parser

func (l *lexer) MatchGrantToken() bool {
  return l.Match([]byte("grant"), GrantToken)
}
