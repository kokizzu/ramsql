package parser

func (l *lexer) MatchReferencesToken() bool {
  return l.Match([]byte("references"), ReferencesToken)
}
