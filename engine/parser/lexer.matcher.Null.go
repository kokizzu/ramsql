package parser

func (l *lexer) MatchNullToken() bool {
  return l.Match([]byte("null"), NullToken)
}
