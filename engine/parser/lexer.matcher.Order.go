package parser

func (l *lexer) MatchOrderToken() bool {
  return l.Match([]byte("order"), OrderToken)
}
