package parser

func (l *lexer) MatchActionToken() bool {
  return l.Match([]byte("action"), ActionToken)
}
