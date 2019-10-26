package lexer

func (l *Lexer) matchDescToken() bool {
  return l.match([]byte("desc"), DescToken)
}
