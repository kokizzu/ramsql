package lexer

func (l *Lexer) MatchNullToken() bool {
  return l.Match([]byte("null"), NullToken)
}
