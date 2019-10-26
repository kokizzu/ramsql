package lexer

func (l *Lexer) MatchLeftDipleToken() bool {
  return l.MatchSingle('<', LeftDipleToken)
}
