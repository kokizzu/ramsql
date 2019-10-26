package lexer

func (l *Lexer) MatchTableToken() bool {
  return l.Match([]byte("table"), TableToken)
}
