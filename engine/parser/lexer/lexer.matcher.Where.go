package lexer

func (l *Lexer) matchWhereToken() bool {
  return l.match([]byte("where"), WhereToken)
}
