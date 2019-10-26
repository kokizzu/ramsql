package lexer

func (l *Lexer) MatchAscToken() bool {
  return l.Match([]byte("asc"), AscToken)
}
