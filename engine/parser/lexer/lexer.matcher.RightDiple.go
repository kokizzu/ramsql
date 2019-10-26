package lexer

func (l *Lexer) matchRightDipleToken() bool {
  return l.matchSingle('>', RightDipleToken)
}
