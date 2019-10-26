package lexer

func (l *Lexer) MatchInsertToken() bool {
  return l.Match([]byte("insert"), InsertToken)
}
