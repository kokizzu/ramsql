package parser

func (l *lexer) MatchTableToken() bool {
  return l.Match([]byte("table"), TableToken)
}
