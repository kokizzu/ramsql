package engine

import (
	"testing"

	"github.com/kokizzu/ramsql/engine/log"
	"github.com/kokizzu/ramsql/engine/parser"
)

func TestCreateTable(t *testing.T) {
	log.UseTestLogger(t)

	e := testEngine(t)
	defer e.Stop()

	createTable(e, t)
}

func TestInsertTable(t *testing.T) {
	log.UseTestLogger(t)

	e := testEngine(t)
	defer e.Stop()

	createTable(e, t)

	query := `INSERT INTO user ('last_name', 'first_name', 'email') VALUES ('Roullon', 'Pierre', 'pierre.roullon@gmail.com')`

	i, err := parser.ParseInstruction(query)
	if err != nil {
		t.Fatalf("Cannot parse query %s : %s", query, err)
	}

	err = e.executeQuery(i[0], &TestEngineConn{})
	if err != nil {
		t.Fatalf("Cannot execute query: %s", err)
	}
}

func createTable(e *Engine, t *testing.T) {
	query := `CREATE TABLE user (
      id INT PRIMARY KEY,
	    last_name TEXT,
	    first_name TEXT,
	    email TEXT,
	    birth_date DATE,
	    country TEXT,
	    town TEXT,
	    zip_code TEXT,
      created_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
      modified_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	)`

	i, err := parser.ParseInstruction(query)
	if err != nil {
		t.Fatalf("Cannot parse query %s : %s", query, err)
	}

	err = e.executeQuery(i[0], &TestEngineConn{})
	if err != nil {
		t.Fatalf("Cannot execute query: %s", err)
	}

}

func TestCreateAndInsertTableWithBacktickedKeyword(t *testing.T) {
	log.UseTestLogger(t)

	e := testEngine(t)
	defer e.Stop()

	query := `CREATE TABLE policy (
      id INT PRIMARY KEY AUTO_INCREMENT,
	    ` + "`action`" + ` TEXT,
	    path TEXT
	)`

	i, err := parser.ParseInstruction(query)
	if err != nil {
		t.Fatalf("Cannot parse query %s : %s", query, err)
	}

	err = e.executeQuery(i[0], &TestEngineConn{})
	if err != nil {
		t.Fatalf("Cannot execute query: %s", err)
	}

	query = `INSERT INTO policy (id, ` + "`action`" + `, path) VALUES (0, 'GET', '/v1/customers')`

	i, err = parser.ParseInstruction(query)
	if err != nil {
		t.Fatalf("Cannot parse query %s : %s", query, err)
	}

	err = e.executeQuery(i[0], &TestEngineConn{})
	if err != nil {
		t.Fatalf("Cannot execute query: %s", err)
	}
}

func TestCreateWithAttributeDefaultNull(t *testing.T) {
	log.UseTestLogger(t)

	e := testEngine(t)
	defer e.Stop()

	query := `CREATE TABLE user (
      id INT PRIMARY KEY AUTO_INCREMENT,
	    last_name TEXT,
	    first_name TEXT,
	    email TEXT DEFAULT NULL,
      created_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
      modified_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	)`

	i, err := parser.ParseInstruction(query)
	if err != nil {
		t.Fatalf("Cannot parse query %s : %s", query, err)
	}

	err = e.executeQuery(i[0], &TestEngineConn{})
	if err != nil {
		t.Fatalf("Cannot execute query: %s", err)
	}
}

func TestCreateWithIndexKeyAsc(t *testing.T) {
	log.UseTestLogger(t)

	e := testEngine(t)
	defer e.Stop()

	query := `CREATE TABLE user (
      id INT PRIMARY KEY AUTO_INCREMENT,
      organization_id INT NOT NULL,
      last_name TEXT,
      first_name TEXT,
      email TEXT DEFAULT NULL,
      created_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
      modified_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
      FOREIGN KEY (organization_id)
        REFERENCES organization (organization_id),
      INDEX idx_organization ( organization_id ASC )
	)`

	i, err := parser.ParseInstruction(query)
	if err != nil {
		t.Fatalf("Cannot parse query %s : %s", query, err)
	}

	err = e.executeQuery(i[0], &TestEngineConn{})
	if err != nil {
		t.Fatalf("Cannot execute query: %s", err)
	}
}

func TestCreateWithForeignKeyOnUpdateNoAction(t *testing.T) {
	log.UseTestLogger(t)

	e := testEngine(t)
	defer e.Stop()

	query := `CREATE TABLE user (
      id INT PRIMARY KEY AUTO_INCREMENT,
      organization_id INT NOT NULL,
      last_name TEXT,
      first_name TEXT,
      email TEXT DEFAULT NULL,
      created_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
      modified_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
      FOREIGN KEY (organization_id)
        REFERENCES organization (organization_id) ON DELETE CASCADE ON UPDATE NO ACTION,
      INDEX idx_organization ( organization_id ASC )
	)`

	i, err := parser.ParseInstruction(query)
	if err != nil {
		t.Fatalf("Cannot parse query %s : %s", query, err)
	}

	err = e.executeQuery(i[0], &TestEngineConn{})
	if err != nil {
		t.Fatalf("Cannot execute query: %s", err)
	}
}
