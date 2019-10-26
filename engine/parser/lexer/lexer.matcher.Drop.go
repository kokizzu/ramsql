package lexer

func (l *Lexer) MatchDropToken() bool {
  return l.Match([]byte("drop"), DropToken)
}
