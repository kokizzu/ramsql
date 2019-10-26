package lexer

func (l *Lexer) MatchEngineToken() bool {
  return l.Match([]byte("engine"), EngineToken)
}
