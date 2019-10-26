package lexer

func (l *Lexer) matchLeftToken() bool {
  return l.match([]byte("left"), LeftToken)
}
