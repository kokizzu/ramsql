package parser

func (l *lexer) MatchJoinToken() bool {
  return l.Match([]byte("join"), JoinToken)
}
