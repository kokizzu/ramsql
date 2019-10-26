package lexer

func (l *Lexer) MatchForToken() bool {
  return l.Match([]byte("for"), ForToken)
}
