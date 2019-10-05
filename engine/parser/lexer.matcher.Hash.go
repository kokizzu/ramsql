package parser

func (l *lexer) MatchHashToken() bool {
  return l.Match([]byte("hash"), HashToken)
}
