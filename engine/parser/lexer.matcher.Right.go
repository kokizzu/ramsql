package parser

func (l *lexer) MatchRightToken() bool {
  return l.Match([]byte("right"), RightToken)
}
