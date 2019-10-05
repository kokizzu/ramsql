package parser

func (l *lexer) MatchNoToken() bool {
  return l.Match([]byte("no"), NoToken)
}
