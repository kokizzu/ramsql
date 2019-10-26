package lexer

func (l *Lexer) MatchValuesToken() bool {
  return l.Match([]byte("values"), ValuesToken)
}
