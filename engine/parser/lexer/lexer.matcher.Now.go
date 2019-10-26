package lexer

func (l *Lexer) matchNowToken() bool {
  return l.match([]byte("now()"), NowToken)
}
