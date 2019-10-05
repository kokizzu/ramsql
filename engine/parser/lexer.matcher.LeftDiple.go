package parser

func (l *lexer) MatchLeftDipleToken() bool {
  return l.MatchSingle('<', LeftDipleToken)
}
