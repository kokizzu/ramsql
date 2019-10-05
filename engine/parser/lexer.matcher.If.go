package parser

func (l *lexer) MatchIfToken() bool {
  return l.Match([]byte("if"), IfToken)
}
