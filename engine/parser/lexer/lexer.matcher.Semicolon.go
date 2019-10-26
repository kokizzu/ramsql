package lexer

func (l *Lexer) MatchSemicolonToken() bool {
  return l.MatchSingle(';', SemicolonToken)
}
