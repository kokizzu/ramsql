package lexer

func (l *Lexer) MatchAndToken() bool {
  return l.Match([]byte("and"), AndToken)
}
