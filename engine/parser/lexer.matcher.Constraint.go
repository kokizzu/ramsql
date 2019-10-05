package parser

func (l *lexer) MatchConstraintToken() bool {
  return l.Match([]byte("constraint"), ConstraintToken)
}
