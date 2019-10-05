package parser

func (l *lexer) MatchEqualityToken() bool {
  return l.MatchSingle('=', EqualityToken)
}
