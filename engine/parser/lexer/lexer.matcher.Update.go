package lexer

func (l *Lexer) MatchUpdateToken() bool {
  return l.Match([]byte("update"), UpdateToken)
}
