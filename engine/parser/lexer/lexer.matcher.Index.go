package lexer

func (l *Lexer) MatchIndexToken() bool {
  return l.Match([]byte("index"), IndexToken)
}
