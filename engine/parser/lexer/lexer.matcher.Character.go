package lexer

func (l *Lexer) MatchCharacterToken() bool {
  return l.Match([]byte("character"), CharacterToken)
}
