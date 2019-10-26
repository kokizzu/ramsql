package lexer

func (l *Lexer) matchIntoToken() bool {
  return l.match([]byte("into"), IntoToken)
}
