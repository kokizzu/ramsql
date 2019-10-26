package lexer

func (l *Lexer) MatchUsingToken() bool {
  return l.Match([]byte("using"), UsingToken)
}
