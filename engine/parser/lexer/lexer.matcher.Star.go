package lexer

func (l *Lexer) matchStarToken() bool {
  return l.matchSingle('*', StarToken)
}
