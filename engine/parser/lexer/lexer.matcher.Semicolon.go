package lexer

func (l *Lexer) matchSemicolonToken() bool {
  return l.matchSingle(';', SemicolonToken)
}
