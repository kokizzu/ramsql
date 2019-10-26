package lexer

func (l *Lexer) MatchJoinToken() bool {
  return l.Match([]byte("join"), JoinToken)
}
