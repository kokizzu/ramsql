package parser

func (l *lexer) MatchLimitToken() bool {
  return l.Match([]byte("limit"), LimitToken)
}
