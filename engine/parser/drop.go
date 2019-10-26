package parser

import (
	"github.com/mlhoyt/ramsql/engine/log"
	"github.com/mlhoyt/ramsql/engine/parser/lexer"
)

func (p *parser) parseDrop() (*Instruction, error) {
	i := &Instruction{}

	// Required: DROP
	dropDecl, err := p.consumeToken(lexer.DropToken)
	if err != nil {
		log.Debug("WTF\n")
		return nil, err
	}
	i.Decls = append(i.Decls, dropDecl)

	// Required: TABLE
	tableDecl, err := p.consumeToken(lexer.TableToken)
	if err != nil {
		log.Debug("Consume table !\n")
		return nil, err
	}
	dropDecl.Add(tableDecl)

	// Optional: IF EXISTS
	if p.is(lexer.IfToken) {
		ifDecl, err := p.consumeToken(lexer.IfToken)
		if err != nil {
			return nil, err
		}
		tableDecl.Add(ifDecl)

		// Required: EXISTS
		if !p.is(lexer.ExistsToken) {
			return nil, p.syntaxError()
		}

		existsDecl, err := p.consumeToken(lexer.ExistsToken)
		if err != nil {
			return nil, err
		}
		ifDecl.Add(existsDecl)
	}

	// Required: <TABLE-NAME>
	nameDecl, err := p.parseQuotedToken()
	if err != nil {
		log.Debug("UH ?\n")
		return nil, err
	}
	tableDecl.Add(nameDecl)

	return i, nil
}
