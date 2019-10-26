package lexer

func (l *Lexer) MatchRightDipleToken() bool {
  return l.MatchSingle('>', RightDipleToken)
}
