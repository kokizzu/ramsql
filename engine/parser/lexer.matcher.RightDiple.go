package parser

func (l *lexer) MatchRightDipleToken() bool {
  return l.MatchSingle('>', RightDipleToken)
}
