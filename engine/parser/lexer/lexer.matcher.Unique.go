package lexer

func (l *Lexer) matchUniqueToken() bool {
  return l.match([]byte("unique"), UniqueToken)
}
