package lexer

func (l *Lexer) matchCommaToken() bool {
  return l.matchSingle(',', CommaToken)
}
