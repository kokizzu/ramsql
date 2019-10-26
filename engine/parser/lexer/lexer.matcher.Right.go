package lexer

func (l *Lexer) matchRightToken() bool {
  return l.match([]byte("right"), RightToken)
}
