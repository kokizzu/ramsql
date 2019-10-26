package lexer

func (l *Lexer) MatchInnerToken() bool {
  return l.Match([]byte("inner"), InnerToken)
}
