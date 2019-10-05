package parser

func (l *lexer) MatchIndexToken() bool {
  return l.Match([]byte("index"), IndexToken)
}
