package parser

func (l *lexer) MatchDeleteToken() bool {
  return l.Match([]byte("delete"), DeleteToken)
}
