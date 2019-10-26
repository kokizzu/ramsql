package lexer

func (l *Lexer) MatchSetToken() bool {
  return l.Match([]byte("set"), SetToken)
}
