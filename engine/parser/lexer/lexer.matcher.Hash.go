package lexer

func (l *Lexer) matchHashToken() bool {
  return l.match([]byte("hash"), HashToken)
}
