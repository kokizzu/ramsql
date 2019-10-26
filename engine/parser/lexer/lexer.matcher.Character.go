package lexer

func (l *Lexer) matchCharacterToken() bool {
  return l.match([]byte("character"), CharacterToken)
}
