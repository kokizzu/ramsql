package lexer

func (l *Lexer) MatchNotToken() bool {
  return l.Match([]byte("not"), NotToken)
}
