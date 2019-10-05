package parser

func (l *lexer) MatchFullToken() bool {
  return l.Match([]byte("full"), FullToken)
}
