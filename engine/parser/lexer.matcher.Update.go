package parser

func (l *lexer) MatchUpdateToken() bool {
  return l.Match([]byte("update"), UpdateToken)
}
