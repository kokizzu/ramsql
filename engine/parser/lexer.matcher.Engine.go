package parser

func (l *lexer) MatchEngineToken() bool {
  return l.Match([]byte("engine"), EngineToken)
}
