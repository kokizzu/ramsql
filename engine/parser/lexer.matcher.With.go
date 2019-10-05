package parser

func (l *lexer) MatchWithToken() bool {
  return l.Match([]byte("with"), WithToken)
}
