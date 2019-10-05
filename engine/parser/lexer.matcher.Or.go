package parser

func (l *lexer) MatchOrToken() bool {
  return l.Match([]byte("or"), OrToken)
}
