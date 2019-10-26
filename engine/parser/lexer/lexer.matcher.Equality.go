package lexer

func (l *Lexer) matchEqualityToken() bool {
  return l.matchSingle('=', EqualityToken)
}
