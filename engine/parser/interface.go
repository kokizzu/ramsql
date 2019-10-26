package parser

import (
	"errors"

	"github.com/mlhoyt/ramsql/engine/parser/lexer"
)

// ParseInstruction calls lexer and parser, then return Decl tree for each instruction
func ParseInstruction(instruction string) ([]Instruction, error) {

	l := lexer.Lexer{}
	tokens, err := l.Lex([]byte(instruction))
	if err != nil {
		return nil, err
	}

	p := NewParser(tokens)
	instructions, err := p.parse()
	if err != nil {
		return nil, err
	}

	if len(instructions) == 0 {
		return nil, errors.New("Error in syntax near " + instruction)
	}

	return instructions, nil
}
