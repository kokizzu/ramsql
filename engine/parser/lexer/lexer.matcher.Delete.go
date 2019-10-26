package lexer

func (l *Lexer) matchDeleteToken() bool {
  return l.match([]byte("delete"), DeleteToken)
}
