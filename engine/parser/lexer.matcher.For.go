package parser

func (l *lexer) MatchForToken() bool {
  return l.Match([]byte("for"), ForToken)
}
