package lexer

import (
	"testing"
	"time"
)

func TestLexerSimple(t *testing.T) {
	query := `CREATE TABLE ` + "`" + `account` + "`" + ``

	lexer := Lexer{}
	decls, err := lexer.Lex([]byte(query))
	if err != nil {
		t.Fatalf("Cannot lex <%s> string", query)
	}

	if len(decls) != 7 {
		t.Fatalf("Lexing failed, expected 7 tokens, got %d", len(decls))
	}
}

func TestParseDate(t *testing.T) {
	data := `2015-09-10T14:03:09.444695269Z`

	_, err := time.Parse(time.RFC3339Nano, data)
	if err != nil {
		t.Fatalf("Cannot parse %s: %s", data, err)
	}
}

func TestLexerDate(t *testing.T) {
	query := `2015-09-10T14:03:09.444695269Z`

	lexer := Lexer{}
	decls, err := lexer.Lex([]byte(query))
	if err != nil {
		t.Fatalf("Cannot lex <%s> string", query)
	}

	if len(decls) != 1 {
		t.Fatalf("Lexing failed, expected 1 tokens, got %d", len(decls))
	}

	if decls[0].Token != StringToken {
		t.Fatalf("Lexing failed, expected a String (%d) tokens, got %d", DateToken, decls[0].Token)
	}
}

func TestLexerWithGTOEandLTOEOperator(t *testing.T) {
	query := `SELECT FROM foo WHERE 1 >= 1 AND 2 <= 3`

	lexer := Lexer{}
	decls, err := lexer.Lex([]byte(query))
	if err != nil {
		t.Fatalf("Cannot lex <%s> string", query)
	}

	if len(decls) != 21 {
		t.Fatalf("Lexing failed, expected 21 tokens, got %d", len(decls))
	}
}

func TestLexerAutoIncrement(t *testing.T) {
	query := `AUTOINCREMENT`

	lexer := Lexer{}
	decls, err := lexer.Lex([]byte(query))
	if err != nil {
		t.Fatalf("Cannot lex <%s> string", query)
	}

	nrExpectedDecls := 1
	if len(decls) != nrExpectedDecls {
		t.Fatalf("Lexing failed, expected %d tokens, got %d", nrExpectedDecls, len(decls))
	}
}

func TestLexerOtherAutoIncrement(t *testing.T) {
	query := `AUTO_INCREMENT`

	lexer := Lexer{}
	decls, err := lexer.Lex([]byte(query))
	if err != nil {
		t.Fatalf("Cannot lex <%s> string", query)
	}

	nrExpectedDecls := 1
	if len(decls) != nrExpectedDecls {
		t.Fatalf("Lexing failed, expected %d tokens, got %d", nrExpectedDecls, len(decls))
	}
}
