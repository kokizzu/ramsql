package parser

func (l *lexer) MatchPrimaryToken() bool {
  return l.Match([]byte("primary"), PrimaryToken)
}
