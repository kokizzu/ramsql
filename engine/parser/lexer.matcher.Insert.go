package parser

func (l *lexer) MatchInsertToken() bool {
  return l.Match([]byte("insert"), InsertToken)
}
