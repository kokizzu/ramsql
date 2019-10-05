package parser

func (l *lexer) MatchAscToken() bool {
  return l.Match([]byte("asc"), AscToken)
}
