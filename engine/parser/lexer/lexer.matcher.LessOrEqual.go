package lexer

func (l *Lexer) matchLessOrEqualToken() bool {
  return l.match([]byte("<="), LessOrEqualToken)
}
