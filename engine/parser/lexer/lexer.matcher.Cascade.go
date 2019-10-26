package lexer

func (l *Lexer) MatchCascadeToken() bool {
  return l.Match([]byte("cascade"), CascadeToken)
}
