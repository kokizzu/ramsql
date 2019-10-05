package parser

func (l *lexer) MatchSemicolonToken() bool {
  return l.MatchSingle(';', SemicolonToken)
}
