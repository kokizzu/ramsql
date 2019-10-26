package lexer

func (l *Lexer) matchIndexToken() bool {
  return l.match([]byte("index"), IndexToken)
}
