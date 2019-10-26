package lexer

func (l *Lexer) matchRestrictToken() bool {
  return l.match([]byte("restrict"), RestrictToken)
}
