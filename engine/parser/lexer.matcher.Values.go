package parser

func (l *lexer) MatchValuesToken() bool {
  return l.Match([]byte("values"), ValuesToken)
}
