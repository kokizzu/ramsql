package lexer

func (l *Lexer) MatchStarToken() bool {
  return l.MatchSingle('*', StarToken)
}
