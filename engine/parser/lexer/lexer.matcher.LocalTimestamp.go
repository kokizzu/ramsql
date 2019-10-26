package lexer

func (l *Lexer) matchLocalTimestampToken() bool {
  return l.match([]byte("localtimestamp"), LocalTimestampToken) ||
     l.match([]byte("current_timestamp"), LocalTimestampToken)
}
