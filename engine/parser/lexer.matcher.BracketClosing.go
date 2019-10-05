package parser

func (l *lexer) MatchBracketClosingToken() bool {
  return l.MatchSingle(')', BracketClosingToken)
}
