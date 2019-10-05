package parser

func (l *lexer) MatchCommaToken() bool {
  return l.MatchSingle(',', CommaToken)
}
