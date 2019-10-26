package lexer

func (l *Lexer) MatchIntoToken() bool {
  return l.Match([]byte("into"), IntoToken)
}
