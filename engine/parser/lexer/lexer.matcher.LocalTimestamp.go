package lexer

func (l *Lexer) MatchLocalTimestampToken() bool {
  return l.Match([]byte("localtimestamp"), LocalTimestampToken) ||
     l.Match([]byte("current_timestamp"), LocalTimestampToken)
}
