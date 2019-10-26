package lexer

func (l *Lexer) MatchPartialToken() bool {
  return l.Match([]byte("partial"), PartialToken)
}
