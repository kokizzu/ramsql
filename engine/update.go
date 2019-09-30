package engine

import (
	"fmt"
	"strings"
	"time"

	"github.com/mlhoyt/ramsql/engine/log"
	"github.com/mlhoyt/ramsql/engine/parser"
	"github.com/mlhoyt/ramsql/engine/protocol"
)

/*
|-> update
	|-> account
	|-> set
	      |-> email
					|-> =
					|-> roger@gmail.com
  |-> where
        |-> id
					|-> =
					|-> 2
*/
func updateExecutor(e *Engine, updateDecl *parser.Decl, conn protocol.EngineConn) error {
	var num int64

	updateDecl.Stringy(0)

	// Fetch table from name and write lock it
	r := e.relation(updateDecl.Decl[0].Lexeme)
	if r == nil {
		return fmt.Errorf("Table %s does not exists", updateDecl.Decl[0].Lexeme)
	}
	r.Lock()
	r.Unlock()

	// Set decl
	values, err := setExecutor(updateDecl.Decl[1])
	if err != nil {
		return err
	}

	// Where decl
	predicates, err := whereExecutor(updateDecl.Decl[2], r.table.name)
	if err != nil {
		return err
	}

	var ok, res bool
	for i := range r.rows {
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
			num++
			err = updateValues(r, i, values)
			if err != nil {
				return err
			}
		}
	}

	return conn.WriteResult(0, num)
}

/*
	|-> set
	      |-> email
					|-> =
					|-> roger@gmail.com
*/
func setExecutor(setDecl *parser.Decl) (map[string]interface{}, error) {

	values := make(map[string]interface{})

	for _, attr := range setDecl.Decl {
		values[attr.Lexeme] = attr.Decl[1].Lexeme
	}

	return values, nil
}

func updateValues(r *Relation, row int, values map[string]interface{}) error {
	for i := range r.table.attributes {
		val, ok := values[r.table.attributes[i].name]
		if !ok {
			switch onUpdateVal := r.table.attributes[i].onUpdateValue.(type) {
			case func() interface{}:
				val = (func() interface{})(onUpdateVal)()
			default:
				continue
			}
		} else {
			log.Debug("Type of '%s' is '%s'\n", r.table.attributes[i].name, r.table.attributes[i].typeName)
			switch strings.ToLower(r.table.attributes[i].typeName) {
			case "timestamp", "localtimestamp":
				switch valVal := val.(type) {
				case func() interface{}:
					val = (func() interface{})(valVal)()
				case time.Time:
					// format time.Time into parsable string
					val = valVal.Format(parser.DateLongFormat)
				case string:
					if valVal == "current_timestamp" || valVal == "now()" {
						val = time.Now().Format(parser.DateLongFormat)
					}
				}
			}
		}
		r.rows[row].Values[i] = fmt.Sprintf("%v", val)
	}

	return nil
}
