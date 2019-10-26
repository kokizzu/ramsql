package lexer

func (l *Lexer) MatchOrderToken() bool {
  return l.Match([]byte("order"), OrderToken)
}
