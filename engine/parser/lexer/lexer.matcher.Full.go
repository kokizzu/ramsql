package lexer

func (l *Lexer) MatchFullToken() bool {
  return l.Match([]byte("full"), FullToken)
}
