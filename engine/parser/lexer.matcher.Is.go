package parser

func (l *lexer) MatchIsToken() bool {
  return l.Match([]byte("is"), IsToken)
}
