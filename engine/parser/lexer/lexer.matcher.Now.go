package lexer

func (l *Lexer) MatchNowToken() bool {
  return l.Match([]byte("now()"), NowToken)
}
