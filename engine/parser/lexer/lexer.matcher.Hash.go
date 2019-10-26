package lexer

func (l *Lexer) MatchHashToken() bool {
  return l.Match([]byte("hash"), HashToken)
}
