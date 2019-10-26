package lexer

func (l *Lexer) MatchLeftToken() bool {
  return l.Match([]byte("left"), LeftToken)
}
