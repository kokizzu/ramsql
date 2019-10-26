package lexer

func (l *Lexer) MatchKeyToken() bool {
  return l.Match([]byte("key"), KeyToken)
}
