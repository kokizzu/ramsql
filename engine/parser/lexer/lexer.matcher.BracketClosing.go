package lexer

func (l *Lexer) matchBracketClosingToken() bool {
  return l.matchSingle(')', BracketClosingToken)
}
