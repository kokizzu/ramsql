package lexer

func (l *Lexer) MatchZoneToken() bool {
  return l.Match([]byte("zone"), ZoneToken)
}
