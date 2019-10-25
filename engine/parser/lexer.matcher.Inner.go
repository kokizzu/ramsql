package parser

func (l *lexer) MatchInnerToken() bool {
  return l.Match([]byte("inner"), InnerToken)
}
