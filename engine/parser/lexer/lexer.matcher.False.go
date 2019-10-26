package lexer

func (l *Lexer) matchFalseToken() bool {
  return l.match([]byte("false"), FalseToken)
}
