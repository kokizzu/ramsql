package parser

func (l *lexer) MatchByToken() bool {
  return l.Match([]byte("by"), ByToken)
}
