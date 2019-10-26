package lexer

func (l *Lexer) MatchGrantToken() bool {
  return l.Match([]byte("grant"), GrantToken)
}
