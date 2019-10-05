package parser

func (l *lexer) MatchForeignToken() bool {
  return l.Match([]byte("foreign"), ForeignToken)
}
