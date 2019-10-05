package parser

func (l *lexer) MatchRestrictToken() bool {
  return l.Match([]byte("restrict"), RestrictToken)
}
