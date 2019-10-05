package parser

func (l *lexer) MatchNowToken() bool {
  return l.Match([]byte("now()"), NowToken)
}
