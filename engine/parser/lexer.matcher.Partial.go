package parser

func (l *lexer) MatchPartialToken() bool {
  return l.Match([]byte("partial"), PartialToken)
}
