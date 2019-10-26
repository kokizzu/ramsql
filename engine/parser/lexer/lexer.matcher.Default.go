package lexer

func (l *Lexer) MatchDefaultToken() bool {
  return l.Match([]byte("default"), DefaultToken)
}
