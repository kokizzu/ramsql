package lexer

func (l *Lexer) MatchInToken() bool {
  return l.Match([]byte("in"), InToken)
}
