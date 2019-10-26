package lexer

func (l *Lexer) matchBacktickToken() bool {
  return l.matchSingle('`', BacktickToken)
}
