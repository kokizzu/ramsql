package engine

import (
	"github.com/kokizzu/ramsql/engine/log"
	"github.com/kokizzu/ramsql/engine/parser"
	"github.com/kokizzu/ramsql/engine/protocol"
)

func truncateExecutor(e *Engine, trDecl *parser.Decl, conn protocol.EngineConn) error {
	log.Debug("truncateExecutor")

	// get tables to be deleted
	table := NewTable(trDecl.Decl[0].Lexeme)

	return truncateTable(e, table, conn)
}
