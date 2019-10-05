package parser

func (l *lexer) MatchFalseToken() bool {
  return l.Match([]byte("false"), FalseToken)
}
