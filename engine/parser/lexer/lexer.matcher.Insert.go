package lexer

func (l *Lexer) matchInsertToken() bool {
  return l.match([]byte("insert"), InsertToken)
}
