package parser

func (l *lexer) MatchBtreeToken() bool {
  return l.Match([]byte("btree"), BtreeToken)
}
