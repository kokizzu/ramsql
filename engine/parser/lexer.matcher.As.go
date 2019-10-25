package parser

func (l *lexer) MatchAsToken() bool {
  return l.Match([]byte("as"), AsToken)
}
