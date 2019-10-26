package lexer

func (l *Lexer) MatchLessOrEqualToken() bool {
  return l.Match([]byte("<="), LessOrEqualToken)
}
