package lexer

func (l *Lexer) matchFullToken() bool {
  return l.match([]byte("full"), FullToken)
}
