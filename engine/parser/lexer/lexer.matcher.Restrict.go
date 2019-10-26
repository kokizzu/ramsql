package lexer

func (l *Lexer) MatchRestrictToken() bool {
  return l.Match([]byte("restrict"), RestrictToken)
}
