package lexer

func (l *Lexer) MatchDescToken() bool {
  return l.Match([]byte("desc"), DescToken)
}
