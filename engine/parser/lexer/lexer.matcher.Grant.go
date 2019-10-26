package lexer

func (l *Lexer) matchGrantToken() bool {
  return l.match([]byte("grant"), GrantToken)
}
