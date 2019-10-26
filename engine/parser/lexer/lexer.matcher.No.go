package lexer

func (l *Lexer) MatchNoToken() bool {
  return l.Match([]byte("no"), NoToken)
}
