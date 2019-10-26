package lexer

func (l *Lexer) MatchRightToken() bool {
  return l.Match([]byte("right"), RightToken)
}
