package lexer

func (l *Lexer) MatchConstraintToken() bool {
  return l.Match([]byte("constraint"), ConstraintToken)
}
