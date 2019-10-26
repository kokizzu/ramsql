package lexer

func (l *Lexer) MatchIsToken() bool {
  return l.Match([]byte("is"), IsToken)
}
