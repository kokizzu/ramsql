package lexer

func (l *Lexer) matchIfToken() bool {
  return l.match([]byte("if"), IfToken)
}
