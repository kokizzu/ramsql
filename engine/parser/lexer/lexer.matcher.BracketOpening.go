package lexer

func (l *Lexer) MatchBracketOpeningToken() bool {
  return l.MatchSingle('(', BracketOpeningToken)
}
