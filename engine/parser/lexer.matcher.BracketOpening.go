package parser

func (l *lexer) MatchBracketOpeningToken() bool {
  return l.MatchSingle('(', BracketOpeningToken)
}
