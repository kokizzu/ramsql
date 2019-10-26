package lexer

func (l *Lexer) matchZoneToken() bool {
  return l.match([]byte("zone"), ZoneToken)
}
