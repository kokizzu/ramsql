package parser

func (l *lexer) MatchWhereToken() bool {
  return l.Match([]byte("where"), WhereToken)
}
