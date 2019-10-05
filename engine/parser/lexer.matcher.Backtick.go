package parser

func (l *lexer) MatchBacktickToken() bool {
  return l.MatchSingle('`', BacktickToken)
}
