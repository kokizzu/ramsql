package lexer

func (l *Lexer) MatchFalseToken() bool {
  return l.Match([]byte("false"), FalseToken)
}
