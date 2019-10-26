package lexer

func (l *Lexer) matchBtreeToken() bool {
  return l.match([]byte("btree"), BtreeToken)
}
