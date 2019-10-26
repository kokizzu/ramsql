package lexer

func (l *Lexer) MatchAsToken() bool {
  return l.Match([]byte("as"), AsToken)
}
