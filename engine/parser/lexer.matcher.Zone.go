package parser

func (l *lexer) MatchZoneToken() bool {
  return l.Match([]byte("zone"), ZoneToken)
}
