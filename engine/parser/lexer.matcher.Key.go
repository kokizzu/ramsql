package parser

func (l *lexer) MatchKeyToken() bool {
  return l.Match([]byte("key"), KeyToken)
}
