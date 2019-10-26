package lexer

func (l *Lexer) MatchDeleteToken() bool {
  return l.Match([]byte("delete"), DeleteToken)
}
