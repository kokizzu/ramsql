package engine

import (
	"fmt"

	"github.com/kokizzu/ramsql/engine/parser"
	"github.com/kokizzu/ramsql/engine/protocol"
)

func ifExecutor(e *Engine, ifDecl *parser.Decl, conn protocol.EngineConn) error {

	if len(ifDecl.Decl) == 0 {
		return fmt.Errorf("malformed condition")
	}

	if e.opsExecutors[ifDecl.Decl[0].Token] != nil {
		return e.opsExecutors[ifDecl.Decl[0].Token](e, ifDecl.Decl[0], conn)
	}

	return fmt.Errorf("error near %v, unknown keyword", ifDecl.Decl[0].Lexeme)
}
