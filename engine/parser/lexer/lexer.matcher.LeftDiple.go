package lexer

func (l *Lexer) matchLeftDipleToken() bool {
  return l.matchSingle('<', LeftDipleToken)
}
