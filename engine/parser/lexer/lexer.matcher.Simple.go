package lexer

func (l *Lexer) MatchSimpleToken() bool {
  return l.Match([]byte("simple"), SimpleToken)
}
