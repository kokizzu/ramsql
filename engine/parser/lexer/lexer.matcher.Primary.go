package lexer

func (l *Lexer) MatchPrimaryToken() bool {
  return l.Match([]byte("primary"), PrimaryToken)
}
