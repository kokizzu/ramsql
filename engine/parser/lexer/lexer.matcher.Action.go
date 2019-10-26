package lexer

func (l *Lexer) MatchActionToken() bool {
  return l.Match([]byte("action"), ActionToken)
}
