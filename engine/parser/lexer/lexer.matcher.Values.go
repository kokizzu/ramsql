package lexer

func (l *Lexer) matchValuesToken() bool {
  return l.match([]byte("values"), ValuesToken)
}
