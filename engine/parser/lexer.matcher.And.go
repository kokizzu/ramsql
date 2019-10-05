package parser

func (l *lexer) MatchAndToken() bool {
  return l.Match([]byte("and"), AndToken)
}
