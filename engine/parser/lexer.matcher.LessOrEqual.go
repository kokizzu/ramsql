package parser

func (l *lexer) MatchLessOrEqualToken() bool {
  return l.Match([]byte("<="), LessOrEqualToken)
}
