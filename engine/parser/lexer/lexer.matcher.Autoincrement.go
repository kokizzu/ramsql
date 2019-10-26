package lexer

func (l *Lexer) MatchAutoincrementToken() bool {
  return l.Match([]byte("autoincrement"), AutoincrementToken) ||
     l.Match([]byte("auto_increment"), AutoincrementToken)
}
