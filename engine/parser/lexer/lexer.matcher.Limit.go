package lexer

func (l *Lexer) MatchLimitToken() bool {
  return l.Match([]byte("limit"), LimitToken)
}
