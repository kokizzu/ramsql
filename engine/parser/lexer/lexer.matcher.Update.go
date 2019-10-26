package lexer

func (l *Lexer) matchUpdateToken() bool {
  return l.match([]byte("update"), UpdateToken)
}
