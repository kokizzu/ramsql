package lexer

func (l *Lexer) matchForeignToken() bool {
  return l.match([]byte("foreign"), ForeignToken)
}
