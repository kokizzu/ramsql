package parser

func (l *lexer) MatchSimpleToken() bool {
  return l.Match([]byte("simple"), SimpleToken)
}
