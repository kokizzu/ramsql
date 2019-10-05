package parser

func (l *lexer) MatchCharacterToken() bool {
  return l.Match([]byte("character"), CharacterToken)
}
