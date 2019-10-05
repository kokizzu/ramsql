package parser

func (l *lexer) MatchStarToken() bool {
  return l.MatchSingle('*', StarToken)
}
