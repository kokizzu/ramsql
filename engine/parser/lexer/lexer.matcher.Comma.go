package lexer

func (l *Lexer) MatchCommaToken() bool {
  return l.MatchSingle(',', CommaToken)
}
