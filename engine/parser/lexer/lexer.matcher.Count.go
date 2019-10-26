package lexer

func (l *Lexer) MatchCountToken() bool {
  return l.Match([]byte("count"), CountToken)
}
