package parser

func (l *lexer) MatchPeriodToken() bool {
  return l.MatchSingle('.', PeriodToken)
}
