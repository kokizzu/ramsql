package lexer

func (l *Lexer) MatchBracketClosingToken() bool {
  return l.MatchSingle(')', BracketClosingToken)
}
