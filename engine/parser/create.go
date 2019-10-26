package parser

import (
	"fmt"
	"strings"

	"github.com/mlhoyt/ramsql/engine/parser/lexer"
)

func (p *Parser) parseCreate() (*Instruction, error) {
	i := &Instruction{}

	// Set CREATE decl
	createDecl := NewDecl(p.cur())
	i.Decls = append(i.Decls, createDecl)

	// After create token, should be either
	// TABLE
	// INDEX
	// ...
	if !p.hasNext() {
		return nil, fmt.Errorf("CREATE token must be followed by TABLE, INDEX")
	}
	p.index++

	switch p.cur().Token {
	case lexer.TableToken:
		d, err := p.parseTable()
		if err != nil {
			return nil, err
		}
		createDecl.Add(d)
		break
	default:
		return nil, fmt.Errorf("Parsing error near <%s>", p.cur().Lexeme)
	}

	return i, nil
}

func (p *Parser) parseTable() (*Decl, error) {
	var err error
	tableDecl := NewDecl(p.cur())
	p.index++

	// Optional: IF NOT EXISTS
	if p.is(lexer.IfToken) {
		ifDecl, err := p.consumeToken(lexer.IfToken)
		if err != nil {
			return nil, err
		}
		tableDecl.Add(ifDecl)

		if p.is(lexer.NotToken) {
			notDecl, err := p.consumeToken(lexer.NotToken)
			if err != nil {
				return nil, err
			}
			ifDecl.Add(notDecl)
			if !p.is(lexer.ExistsToken) {
				return nil, p.syntaxError()
			}
			existsDecl, err := p.consumeToken(lexer.ExistsToken)
			if err != nil {
				return nil, err
			}
			notDecl.Add(existsDecl)
		}
	}

	// Required: <TABLE-NAME>
	nameTable, err := p.parseAttribute()
	if err != nil {
		return nil, p.syntaxError()
	}
	tableDecl.Add(nameTable)

	// Required: '(' (Opening Parenthesis)
	if !p.hasNext() || p.cur().Token != lexer.BracketOpeningToken {
		return nil, fmt.Errorf("Table name token must be followed by table definition")
	}
	p.index++

	// Required: <TABLE-BODY>
	for p.index < p.tokenLen {

		// ')' (Closing parenthesis)
		if p.cur().Token == lexer.BracketClosingToken {
			p.consumeToken(lexer.BracketClosingToken)
			break
		}

		// (CONSTRAINT <CONSTRAINT-NAME>?)? ...
		if p.cur().Token == lexer.ConstraintToken {
			_, err := p.parseTableConstraint()
			if err != nil {
				return nil, err
			}

			// PRIMARY KEY ( <INDEX-KEY> [, ...] )
		} else if p.cur().Token == lexer.PrimaryToken {
			_, err := p.parsePrimaryKey()
			if err != nil {
				return nil, err
			}

			// UNIQUE [INDEX | KEY] ...
		} else if p.cur().Token == lexer.UniqueToken {
			_, err := p.consumeToken(lexer.UniqueToken)
			if err != nil {
				return nil, err
			}

			_, err = p.parseTableIndex()
			if err != nil {
				return nil, err
			}

			// { INDEX | KEY } [ index_name ] [?:index_type USING { BTREE | HASH } ] '(' { col_name [ '(' length ')' ] | '(' expr ')' } [ ASC | DESC ] ',' ... ')' [?:index_option ... ]
		} else if p.cur().Token == lexer.IndexToken || p.cur().Token == lexer.KeyToken {
			_, err := p.parseTableIndex()
			if err != nil {
				return nil, err
			}

			// FOREIGN KEY ...
		} else if p.cur().Token == lexer.ForeignToken {
			_, err := p.parseTableForeignKey()
			if err != nil {
				return nil, err
			}

			// <TABLE-ATTRIBUTE>
		} else {
			// New attribute name
			newAttribute, err := p.parseQuotedToken()
			if err != nil {
				return nil, err
			}
			tableDecl.Add(newAttribute)

			newAttributeType, err := p.parseType()
			if err != nil {
				return nil, err
			}
			newAttribute.Add(newAttributeType)

			// All the following tokens until bracket or comma are column constraints.
			// Column constraints can be listed in any order.
			for p.isNot(lexer.BracketClosingToken, lexer.CommaToken) {
				switch p.cur().Token {
				case lexer.UniqueToken: // UNIQUE
					uniqueDecl, err := p.consumeToken(lexer.UniqueToken)
					if err != nil {
						return nil, err
					}
					newAttribute.Add(uniqueDecl)
				case lexer.NotToken: // NOT NULL
					if _, err = p.isNext(lexer.NullToken); err == nil {
						notDecl, err := p.consumeToken(lexer.NotToken)
						if err != nil {
							return nil, err
						}
						newAttribute.Add(notDecl)
						nullDecl, err := p.consumeToken(lexer.NullToken)
						if err != nil {
							return nil, err
						}
						notDecl.Add(nullDecl)
					}
				case lexer.NullToken: // NULL
					nullDecl, err := p.consumeToken(lexer.NullToken)
					if err != nil {
						return nil, err
					}

					newAttribute.Add(nullDecl)
				case lexer.PrimaryToken: // PRIMARY KEY
					if _, err = p.isNext(lexer.KeyToken); err == nil {
						newPrimary := NewDecl(p.cur())
						newAttribute.Add(newPrimary)

						if err = p.next(); err != nil {
							return nil, fmt.Errorf("Unexpected end")
						}

						newKey := NewDecl(p.cur())
						newPrimary.Add(newKey)

						if err = p.next(); err != nil {
							return nil, fmt.Errorf("Unexpected end")
						}
					}
				case lexer.AutoincrementToken:
					autoincDecl, err := p.consumeToken(lexer.AutoincrementToken)
					if err != nil {
						return nil, err
					}
					newAttribute.Add(autoincDecl)
				case lexer.WithToken: // WITH TIME ZONE
					if strings.ToLower(newAttributeType.Lexeme) == "timestamp" {
						withDecl, err := p.consumeToken(lexer.WithToken)
						if err != nil {
							return nil, err
						}
						timeDecl, err := p.consumeToken(lexer.TimeToken)
						if err != nil {
							return nil, err
						}
						zoneDecl, err := p.consumeToken(lexer.ZoneToken)
						if err != nil {
							return nil, err
						}
						newAttributeType.Add(withDecl)
						withDecl.Add(timeDecl)
						timeDecl.Add(zoneDecl)
					}
				case lexer.DefaultToken: // DEFAULT <VALUE>
					dDecl, err := p.consumeToken(lexer.DefaultToken)
					if err != nil {
						return nil, err
					}
					newAttribute.Add(dDecl)
					vDecl, err := p.consumeToken(lexer.FalseToken, lexer.StringToken, lexer.NumberToken, lexer.LocalTimestampToken)
					if err != nil {
						return nil, err
					}
					dDecl.Add(vDecl)
				case lexer.OnToken: // ON UPDATE <VALUE>
					onDecl, err := p.consumeToken(lexer.OnToken)
					if err != nil {
						return nil, err
					}

					updateDecl, err := p.consumeToken(lexer.UpdateToken)
					if err != nil {
						return nil, err
					}

					vDecl, err := p.consumeToken(lexer.FalseToken, lexer.StringToken, lexer.NumberToken, lexer.LocalTimestampToken)
					if err != nil {
						return nil, err
					}

					onDecl.Add(updateDecl)
					updateDecl.Add(vDecl)
					newAttribute.Add(onDecl)
				default:
					// Unknown column constraint
					return nil, p.syntaxError()
				}
			}
		}

		// Comma means continue to next table column
		// NOTE: With this the parser accepts ", )" and happily proceeds but this is not valid SQL (AFAIK)
		if p.cur().Token == lexer.CommaToken {
			p.index++
		}
	}

	// Optional: <TABLE-OPTIONS> - these can be listed in any order
tableOptions:
	for p.index < p.tokenLen {
		switch p.cur().Token {
		case lexer.EngineToken: // ENGINE [=] value
			engineDecl, err := p.consumeToken(lexer.EngineToken)
			if err != nil {
				return nil, err
			}

			if p.cur().Token == lexer.EqualityToken {
				if err = p.next(); err != nil {
					return nil, err
				}
			}

			vDecl, err := p.consumeToken(lexer.FalseToken, lexer.StringToken, lexer.NumberToken)
			if err != nil {
				return nil, err
			}

			engineDecl.Add(vDecl)
			// TODO: tableDecl.Add(engineDecl)

		case lexer.DefaultToken: // [DEFAULT] (CHARACTER SET, CHARSET, COLLATE) [=] value
			if err := p.next(); err != nil {
				return nil, err
			}

			switch p.cur().Token {
			case lexer.CharsetToken: // CHARSET [=] value
				if err := p.next(); err != nil {
					return nil, err
				}
				charDecl := NewDecl(lexer.Token{Token: lexer.CharacterToken, Lexeme: "character"})
				setDecl := NewDecl(lexer.Token{Token: lexer.SetToken, Lexeme: "set"})

				if p.cur().Token == lexer.EqualityToken {
					if err := p.next(); err != nil {
						return nil, err
					}
				}

				vDecl, err := p.consumeToken(lexer.StringToken)
				if err != nil {
					return nil, err
				}

				charDecl.Add(setDecl)
				setDecl.Add(vDecl)
				// TODO: tableDecl.Add(charDecl)

			case lexer.CharacterToken: // CHARACTER SET [=] value
				charDecl, err := p.consumeToken(lexer.CharacterToken)
				if err != nil {
					return nil, err
				}

				setDecl, err := p.consumeToken(lexer.SetToken)
				if err != nil {
					return nil, err
				}

				if p.cur().Token == lexer.EqualityToken {
					if err := p.next(); err != nil {
						return nil, err
					}
				}

				vDecl, err := p.consumeToken(lexer.StringToken)
				if err != nil {
					return nil, err
				}

				charDecl.Add(setDecl)
				setDecl.Add(vDecl)
				// TODO: tableDecl.Add(charDecl)
			default:
				// Unknown 'table_option'
				return nil, p.syntaxError()
			}

		case lexer.CharsetToken: // CHARSET [=] value
			if err := p.next(); err != nil {
				return nil, err
			}
			charDecl := NewDecl(lexer.Token{Token: lexer.CharacterToken, Lexeme: "character"})
			setDecl := NewDecl(lexer.Token{Token: lexer.SetToken, Lexeme: "set"})

			if p.cur().Token == lexer.EqualityToken {
				if err := p.next(); err != nil {
					return nil, err
				}
			}

			vDecl, err := p.consumeToken(lexer.StringToken)
			if err != nil {
				return nil, err
			}

			charDecl.Add(setDecl)
			setDecl.Add(vDecl)
			// TODO: tableDecl.Add(charDecl)

		case lexer.CharacterToken: // CHARACTER SET [=] value
			charDecl, err := p.consumeToken(lexer.CharacterToken)
			if err != nil {
				return nil, err
			}

			setDecl, err := p.consumeToken(lexer.SetToken)
			if err != nil {
				return nil, err
			}

			if p.cur().Token == lexer.EqualityToken {
				if err := p.next(); err != nil {
					return nil, err
				}
			}

			vDecl, err := p.consumeToken(lexer.StringToken)
			if err != nil {
				return nil, err
			}

			charDecl.Add(setDecl)
			setDecl.Add(vDecl)
			// TODO: tableDecl.Add(charDecl)

		case lexer.SemicolonToken: // semicolon means end of instruction
			// Important NOT to consume the semicolon token

			break tableOptions

		default: // Does not appear to be a 'table_constraint' so stop processing instruction
			break tableOptions
		}
	}

	return tableDecl, nil
}

// parseTableConstraint processes tokens that should define a table constraint
// CONSTRAINT <CONSTRAINT-NAME>? ...
func (p *Parser) parseTableConstraint() (*Decl, error) {
	constraintDecl, err := p.consumeToken(lexer.ConstraintToken)
	if err != nil {
		return nil, err
	}

	// Optional: <CONSTRAINT-NAME>
	if p.is(lexer.StringToken) {
		_, err := p.consumeToken(lexer.StringToken)
		if err != nil {
			return nil, err
		}
	}

	switch p.cur().Token {
	case lexer.PrimaryToken:
		_, err := p.parsePrimaryKey()
		if err != nil {
			return nil, err
		}
	case lexer.UniqueToken:
		_, err := p.consumeToken(lexer.UniqueToken)
		if err != nil {
			return nil, err
		}

		_, err = p.parseTableIndex()
		if err != nil {
			return nil, err
		}
	case lexer.ForeignToken:
		_, err := p.parseTableForeignKey()
		if err != nil {
			return nil, err
		}
	default:
		// Unknown constraint type
		return nil, p.syntaxError()
	}

	return constraintDecl, nil
}

func (p *Parser) parsePrimaryKey() (*Decl, error) {
	primaryDecl, err := p.consumeToken(lexer.PrimaryToken)
	if err != nil {
		return nil, err
	}

	keyDecl, err := p.consumeToken(lexer.KeyToken)
	if err != nil {
		return nil, err
	}
	primaryDecl.Add(keyDecl)

	_, err = p.consumeToken(lexer.BracketOpeningToken)
	if err != nil {
		return nil, err
	}

	for {
		d, err := p.parseQuotedToken()
		if err != nil {
			return nil, err
		}

		d, err = p.consumeToken(lexer.CommaToken, lexer.BracketClosingToken)
		if err != nil {
			return nil, err
		}
		if d.Token == lexer.BracketClosingToken {
			break
		}
	}

	return primaryDecl, nil
}

// parseTableIndex processes tokens that should define a table index
// { INDEX | KEY } [ index_name ] [?:index_type USING { BTREE | HASH } ] '(' { col_name [ '(' length ')' ] | '(' expr ')' } [ ASC | DESC ] ',' ... ')' [?:index_option ... ]
func (p *Parser) parseTableIndex() (*Decl, error) {
	indexDecl := NewDecl(lexer.Token{Token: lexer.IndexToken, Lexeme: "index"})

	// Required: { INDEX | KEY }
	switch p.cur().Token {
	case lexer.IndexToken, lexer.KeyToken:
		_, err := p.consumeToken(lexer.IndexToken, lexer.KeyToken)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("Table INDEX definition must start with INDEX or KEY")
	}

	// Optional: <INDEX-NAME>
	if p.is(lexer.StringToken) {
		_, err := p.consumeToken(lexer.StringToken)
		if err != nil {
			return nil, err
		}
	}

	// Optional: <INDEX-TYPE> := USING { BTREE | HASH }
	if p.is(lexer.UsingToken) {
		_, err := p.consumeToken(lexer.UsingToken)
		if err != nil {
			return nil, err
		}

		if p.is(lexer.BtreeToken) {
			_, err := p.consumeToken(lexer.BtreeToken)
			if err != nil {
				return nil, err
			}
		} else if p.is(lexer.HashToken) {
			_, err := p.consumeToken(lexer.HashToken)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, p.syntaxError()
		}
	}

	// Required: '('
	_, err := p.consumeToken(lexer.BracketOpeningToken)
	if err != nil {
		return nil, err
	}

	// Required: <INDEX-KEY> [, <INDEX-KEY>]* ')'
	for {
		_, err := p.consumeToken(lexer.StringToken)
		if err != nil {
			return nil, err
		}

		n, err := p.consumeToken(lexer.CommaToken, lexer.BracketClosingToken)
		if err != nil {
			return nil, err
		}
		if n.Token == lexer.BracketClosingToken {
			break
		}
	}

	return indexDecl, nil
}

// parseTableForeignKey processes tokens that should define a table foreign key
// FOREIGN KEY ...
func (p *Parser) parseTableForeignKey() (*Decl, error) {
	// Required: FOREIGN
	foreignDecl, err := p.consumeToken(lexer.ForeignToken)
	if err != nil {
		return nil, err
	}

	// Required: KEY
	keyDecl, err := p.consumeToken(lexer.KeyToken)
	if err != nil {
		return nil, err
	}
	foreignDecl.Add(keyDecl)

	// Optional: <FK-NAME>
	if p.is(lexer.StringToken) {
		_, err := p.consumeToken(lexer.StringToken)
		if err != nil {
			return nil, err
		}
	}

	// Required: '('
	_, err = p.consumeToken(lexer.BracketOpeningToken)
	if err != nil {
		return nil, err
	}

	// Required: <FK-INDEX> [, <FK-INDEX>]* ')'
	for {
		_, err := p.consumeToken(lexer.StringToken)
		if err != nil {
			return nil, err
		}

		n, err := p.consumeToken(lexer.CommaToken, lexer.BracketClosingToken)
		if err != nil {
			return nil, err
		}
		if n.Token == lexer.BracketClosingToken {
			break
		}
	}

	// Optional: REFERENCES ...
	if p.is(lexer.ReferencesToken) {
		_, err := p.parseTableReference()
		if err != nil {
			return nil, err
		}
	}

	return foreignDecl, nil
}

// parseTableReference processes tokens that should define a table reference
// REFERENCES ...
func (p *Parser) parseTableReference() (*Decl, error) {
	// Required: REFERENCES
	referencesDecl, err := p.consumeToken(lexer.ReferencesToken)
	if err != nil {
		return nil, err
	}

	// Required: <TABLE-NAME>
	_, err = p.consumeToken(lexer.StringToken)
	if err != nil {
		return nil, err
	}

	// Required: '('
	_, err = p.consumeToken(lexer.BracketOpeningToken)
	if err != nil {
		return nil, err
	}

	// Required: <KEY-PART> [, <KEY-PART>]* ')'
	for {
		_, err := p.consumeToken(lexer.StringToken)
		if err != nil {
			return nil, err
		}

		n, err := p.consumeToken(lexer.CommaToken, lexer.BracketClosingToken)
		if err != nil {
			return nil, err
		}
		if n.Token == lexer.BracketClosingToken {
			break
		}
	}

	// Optional: MATCH ...
	if p.is(lexer.MatchToken) {
		_, err := p.consumeToken(lexer.MatchToken)
		if err != nil {
			return nil, err
		}

		switch p.cur().Token {
		case lexer.FullToken:
			_, err := p.consumeToken(lexer.FullToken)
			if err != nil {
				return nil, err
			}
		case lexer.PartialToken:
			_, err := p.consumeToken(lexer.PartialToken)
			if err != nil {
				return nil, err
			}
		case lexer.SimpleToken:
			_, err := p.consumeToken(lexer.SimpleToken)
			if err != nil {
				return nil, err
			}
		default:
			// Unknown match type
			return nil, p.syntaxError()
		}
	}

	// Optional: ON ...
	if p.is(lexer.OnToken) {
		_, err := p.consumeToken(lexer.OnToken)
		if err != nil {
			return nil, err
		}

		switch p.cur().Token {
		case lexer.UpdateToken:
			_, err := p.consumeToken(lexer.UpdateToken)
			if err != nil {
				return nil, err
			}

			_, err = p.parseTableReferenceOption()
			if err != nil {
				return nil, err
			}
		case lexer.DeleteToken:
			_, err := p.consumeToken(lexer.DeleteToken)
			if err != nil {
				return nil, err
			}

			_, err = p.parseTableReferenceOption()
			if err != nil {
				return nil, err
			}
		default:
			// Unknown on reference option type
			return nil, p.syntaxError()
		}
	}

	return referencesDecl, nil
}

// parseTableReferenceOption processes tokens that should define a table reference option
func (p *Parser) parseTableReferenceOption() (*Decl, error) {
	switch p.cur().Token {
	case lexer.RestrictToken:
		d, err := p.consumeToken(lexer.RestrictToken)
		if err != nil {
			return nil, err
		}
		return d, nil
	case lexer.CascadeToken:
		d, err := p.consumeToken(lexer.CascadeToken)
		if err != nil {
			return nil, err
		}
		return d, nil
	case lexer.SetToken:
		d, err := p.consumeToken(lexer.SetToken)
		if err != nil {
			return nil, err
		}

		switch p.cur().Token {
		case lexer.NullToken:
			n, err := p.consumeToken(lexer.NullToken)
			if err != nil {
				return nil, err
			}
			d.Add(n)
			return d, nil
		case lexer.DefaultToken:
			n, err := p.consumeToken(lexer.DefaultToken)
			if err != nil {
				return nil, err
			}
			d.Add(n)
			return d, nil
		default:
			// Unknown option
			return nil, p.syntaxError()
		}
	case lexer.NoToken:
		d, err := p.consumeToken(lexer.NoToken)
		if err != nil {
			return nil, err
		}

		n, err := p.consumeToken(lexer.ActionToken)
		if err != nil {
			return nil, err
		}
		d.Add(n)
		return d, nil
	default:
		// Unknown option type
		return nil, p.syntaxError()
	}
}
