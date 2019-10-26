package lexer

func (l *Lexer) matchPeriodToken() bool {
  return l.matchSingle('.', PeriodToken)
}
