package lexer

func (l *Lexer) matchAsToken() bool {
  return l.match([]byte("as"), AsToken)
}
