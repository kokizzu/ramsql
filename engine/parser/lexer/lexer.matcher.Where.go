package lexer

func (l *Lexer) MatchWhereToken() bool {
  return l.Match([]byte("where"), WhereToken)
}
