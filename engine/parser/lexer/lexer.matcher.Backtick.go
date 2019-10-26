package lexer

func (l *Lexer) MatchBacktickToken() bool {
  return l.MatchSingle('`', BacktickToken)
}
