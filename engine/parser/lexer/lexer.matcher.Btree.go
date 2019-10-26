package lexer

func (l *Lexer) MatchBtreeToken() bool {
  return l.Match([]byte("btree"), BtreeToken)
}
