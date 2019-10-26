package lexer

func (l *Lexer) matchInnerToken() bool {
  return l.match([]byte("inner"), InnerToken)
}
