package lexer

func (l *Lexer) MatchIfToken() bool {
  return l.Match([]byte("if"), IfToken)
}
