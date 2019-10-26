package lexer

func (l *Lexer) MatchByToken() bool {
  return l.Match([]byte("by"), ByToken)
}
