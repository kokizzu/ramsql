package lexer

func (l *Lexer) matchBracketOpeningToken() bool {
  return l.matchSingle('(', BracketOpeningToken)
}
