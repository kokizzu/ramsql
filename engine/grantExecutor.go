package engine

import (
	"github.com/kokizzu/ramsql/engine/parser"
	"github.com/kokizzu/ramsql/engine/protocol"
)

func grantExecutor(e *Engine, decl *parser.Decl, conn protocol.EngineConn) error {
	return conn.WriteResult(0, 0)
}
