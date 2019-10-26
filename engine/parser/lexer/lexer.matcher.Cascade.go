package lexer

func (l *Lexer) matchCascadeToken() bool {
  return l.match([]byte("cascade"), CascadeToken)
}
