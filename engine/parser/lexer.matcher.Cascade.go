package parser

func (l *lexer) MatchCascadeToken() bool {
  return l.Match([]byte("cascade"), CascadeToken)
}
