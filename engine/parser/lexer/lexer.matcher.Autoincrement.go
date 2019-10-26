package lexer

func (l *Lexer) matchAutoincrementToken() bool {
  return l.match([]byte("autoincrement"), AutoincrementToken) ||
     l.match([]byte("auto_increment"), AutoincrementToken)
}
