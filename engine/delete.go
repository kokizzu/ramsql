package engine

import (
	"fmt"

	"github.com/kokizzu/ramsql/engine/protocol"
)

func deleteRows(e *Engine, tables []*Table, conn protocol.EngineConn, predicates []Predicate) error {
	var rowsDeleted int64

	r := e.relation(tables[0].name)
	if r == nil {
		return fmt.Errorf("Table %s not found", tables[0].name)
	}
	r.Lock()
	defer r.Unlock()

	var ok, res bool
	var err error
	lenRows := len(r.rows)
	for i := 0; i < lenRows; i++ {
		ok = true
		// If the row validate all predicates, write it
		for _, predicate := range predicates {
			if res, err = predicate.Evaluate(r.rows[i], r.table); err != nil {
				return err
			}
			if res == false {
				ok = false
				continue
			}
		}

		if ok {
			switch i {
			case 0:
				r.rows = r.rows[1:]
			case lenRows - 1:
				r.rows = r.rows[:lenRows-1]
			default:
				r.rows = append(r.rows[:i], r.rows[i+1:]...)
				i--
			}
			lenRows--
			rowsDeleted++
		}
	}

	return conn.WriteResult(0, rowsDeleted)
}
