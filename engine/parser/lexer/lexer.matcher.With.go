package lexer

func (l *Lexer) MatchWithToken() bool {
  return l.Match([]byte("with"), WithToken)
}
