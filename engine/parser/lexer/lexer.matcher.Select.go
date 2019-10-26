package lexer

func (l *Lexer) MatchSelectToken() bool {
  return l.Match([]byte("select"), SelectToken)
}
