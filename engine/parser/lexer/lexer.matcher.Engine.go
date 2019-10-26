package lexer

func (l *Lexer) matchEngineToken() bool {
  return l.match([]byte("engine"), EngineToken)
}
