package lexer

func (l *Lexer) MatchOrToken() bool {
  return l.Match([]byte("or"), OrToken)
}
