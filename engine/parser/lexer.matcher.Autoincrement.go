package parser

func (l *lexer) MatchAutoincrementToken() bool {
  return l.Match([]byte("autoincrement"), AutoincrementToken) ||
     l.Match([]byte("auto_increment"), AutoincrementToken)
}
