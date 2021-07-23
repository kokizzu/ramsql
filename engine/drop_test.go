package engine

import (
	"testing"

	"github.com/kokizzu/ramsql/engine/log"
)

func TestDropAfterCreate(t *testing.T) {
	log.UseTestLogger(t)

	e := testEngine(t)
	defer e.Stop()

	err := parseAndExecuteQuery(t, e, "CREATE TABLE account (id INT, email TEXT)")
	if err != nil {
		t.Fatal(err)
	}

	err = parseAndExecuteQuery(t, e, "DROP TABLE account")
	if err != nil {
		t.Fatal(err)
	}
}

func TestDropIfExistsAfterCreate(t *testing.T) {
	log.UseTestLogger(t)

	e := testEngine(t)
	defer e.Stop()

	err := parseAndExecuteQuery(t, e, "CREATE TABLE account (id INT, email TEXT)")
	if err != nil {
		t.Fatal(err)
	}

	err = parseAndExecuteQuery(t, e, "DROP TABLE IF EXISTS account")
	if err != nil {
		t.Fatal(err)
	}
}

func TestDropIfExistsWithoutCreate(t *testing.T) {
	log.UseTestLogger(t)

	e := testEngine(t)
	defer e.Stop()

	err := parseAndExecuteQuery(t, e, "DROP TABLE IF EXISTS account")
	if err != nil {
		t.Fatal(err)
	}
}
