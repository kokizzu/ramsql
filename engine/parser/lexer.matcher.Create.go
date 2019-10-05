package parser

func (l *lexer) MatchCreateToken() bool {
  return l.Match([]byte("create"), CreateToken)
}
