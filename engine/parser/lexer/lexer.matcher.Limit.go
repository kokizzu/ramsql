package lexer

func (l *Lexer) matchLimitToken() bool {
  return l.match([]byte("limit"), LimitToken)
}
