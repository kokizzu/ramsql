package lexer

func (l *Lexer) matchTableToken() bool {
  return l.match([]byte("table"), TableToken)
}
