package engine

import (
	"fmt"
	"testing"

	"github.com/kokizzu/ramsql/engine/parser"
)

func parseAndExecuteQuery(t *testing.T, e *Engine, query string) error {
	i, err := parser.ParseInstruction(query)
	if err != nil {
		return fmt.Errorf("Cannot parse SQL %s : %s", query, err)
	}

	err = e.executeQuery(i[0], &TestEngineConn{})
	if err != nil {
		return fmt.Errorf("Cannot execute SQL: %s", err)
	}

	return nil
}
