package lexer

func (l *Lexer) MatchOnToken() bool {
  return l.Match([]byte("on"), OnToken)
}
