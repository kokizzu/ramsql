package lexer

func (l *Lexer) matchReferencesToken() bool {
  return l.match([]byte("references"), ReferencesToken)
}
