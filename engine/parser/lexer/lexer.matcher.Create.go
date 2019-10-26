package lexer

func (l *Lexer) MatchCreateToken() bool {
  return l.Match([]byte("create"), CreateToken)
}
