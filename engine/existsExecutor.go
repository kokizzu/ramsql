package engine

import (
	"github.com/kokizzu/ramsql/engine/parser"
	"github.com/kokizzu/ramsql/engine/protocol"
)

func existsExecutor(e *Engine, tableDecl *parser.Decl, conn protocol.EngineConn) error {
	return nil
}
