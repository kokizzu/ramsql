package lexer

func (l *Lexer) MatchEqualityToken() bool {
  return l.MatchSingle('=', EqualityToken)
}
