package lexer

func (l *Lexer) matchConstraintToken() bool {
  return l.match([]byte("constraint"), ConstraintToken)
}
