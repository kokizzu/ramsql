package parser

func (l *lexer) MatchDefaultToken() bool {
  return l.Match([]byte("default"), DefaultToken)
}
