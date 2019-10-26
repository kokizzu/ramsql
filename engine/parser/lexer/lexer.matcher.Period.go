package lexer

func (l *Lexer) MatchPeriodToken() bool {
  return l.MatchSingle('.', PeriodToken)
}
