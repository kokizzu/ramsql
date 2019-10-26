package lexer

func (l *Lexer) MatchReferencesToken() bool {
  return l.Match([]byte("references"), ReferencesToken)
}
