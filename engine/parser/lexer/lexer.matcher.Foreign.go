package lexer

func (l *Lexer) MatchForeignToken() bool {
  return l.Match([]byte("foreign"), ForeignToken)
}
