package parser

func (l *lexer) MatchIntoToken() bool {
  return l.Match([]byte("into"), IntoToken)
}
