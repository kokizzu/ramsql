package engine_test

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/kokizzu/ramsql/engine/log"
)

const createTableQuery = `CREATE TABLE user (
  id INT PRIMARY KEY,
  last_name TEXT,
  first_name TEXT,
  email TEXT NULL,
  birth_date DATE NOT NULL,
  country TEXT,
  town TEXT,
  zip_code TEXT,
  created_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  modified_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
)`

func TestSelectNoOp(t *testing.T) {
	log.UseTestLogger(t)
	db, err := sql.Open("ramsql", "TestSelectNoOp")
	if err != nil {
		t.Fatalf("sql.Open : Error : %s\n", err)
	}
	defer db.Close()

	batch := []string{
		`CREATE TABLE account (id BIGSERIAL, email TEXT)`,
		`INSERT INTO account (email) VALUES ("foo@bar.com")`,
		`INSERT INTO account (email) VALUES ("bar@bar.com")`,
		`INSERT INTO account (email) VALUES ("foobar@bar.com")`,
		`INSERT INTO account (email) VALUES ("babar@bar.com")`,
	}

	for _, b := range batch {
		_, err = db.Exec(b)
		if err != nil {
			t.Fatalf("sql.Exec: Error: %s\n", err)
		}
	}

	query := `SELECT * from account WHERE 1 = 1`
	rows, err := db.Query(query)
	if err != nil {
		t.Fatalf("cannot create table: %s", err)
	}

	nb := 0
	for rows.Next() {
		nb++
	}

	if nb != 4 {
		t.Fatalf("Expected 4 rows, got %d", nb)
	}

}

func TestSelect(t *testing.T) {
	log.UseTestLogger(t)
	db, err := sql.Open("ramsql", "TestSelect")
	if err != nil {
		t.Fatalf("sql.Open : Error : %s\n", err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE TABLE account (id INT, email TEXT)")
	if err != nil {
		t.Fatalf("sql.Exec: Error: %s\n", err)
	}

	_, err = db.Exec("INSERT INTO account ('id', 'email') VALUES (2, 'bar@bar.com')")
	if err != nil {
		t.Fatalf("Cannot insert into table account: %s", err)
	}

	_, err = db.Exec("INSERT INTO account ('id', 'email') VALUES (1, 'foo@bar.com')")
	if err != nil {
		t.Fatalf("Cannot insert into table account: %s", err)
	}

	rows, err := db.Query("SELECT * FROM account WHERE email = $1", "foo@bar.com")
	if err != nil {
		t.Fatalf("sql.Query error : %s", err)
	}

	columns, err := rows.Columns()
	if err != nil {
		t.Fatalf("rows.Column : %s", err)
		return
	}

	if len(columns) != 2 {
		t.Fatalf("Expected 2 columns, got %d", len(columns))
	}

	row := db.QueryRow("SELECT * FROM account WHERE email = $1", "foo@bar.com")
	if row == nil {
		t.Fatalf("sql.QueryRow error")
	}

	var email string
	var id int
	err = row.Scan(&id, &email)
	if err != nil {
		t.Fatalf("row.Scan: %s", err)
	}

	if id != 1 {
		t.Fatalf("Expected id = 1, got %d", id)
	}

	if email != "foo@bar.com" {
		t.Fatalf("Expected email = <foo@bar.com>, got <%s>", email)
	}

	err = db.Close()
	if err != nil {
		t.Fatalf("sql.Close : Error : %s\n", err)
	}

}

func TestCreateInsertSelectTableWithBooleanFields(t *testing.T) {
	log.UseTestLogger(t)
	db, err := sql.Open("ramsql", "TestCreateInsertSelectTableWithBooleanFields")
	if err != nil {
		t.Fatalf("sql.Open : Error : %s\n", err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE TABLE account (id INT PRIMARY KEY AUTO_INCREMENT, is_enabled BOOLEAN NOT NULL, is_active BOOLEAN NOT NULL)")
	if err != nil {
		t.Fatalf("Cannot create table account: %s", err)
	}

	_, err = db.Exec("INSERT INTO account (is_enabled, is_active) VALUES (true, false)")
	if err != nil {
		t.Fatalf("Cannot insert into table account: %s", err)
	}

	row := db.QueryRow("SELECT * FROM account WHERE id = 1")
	if row == nil {
		t.Fatalf("Cannot select from table account")
	}
	fmt.Printf("[DEBUG] function:TestCreateInsertSelectTableWithBooleanFields select-row:%v\n", row)

	var id int
	var isEnabled bool
	var isActive bool
	err = row.Scan(&id, &isEnabled, &isActive)
	if err != nil {
		t.Fatalf("row.Scan: %s", err)
	}

	if id != 1 {
		t.Fatalf("Expected id = 1, got %d", id)
	}

	if isEnabled != true {
		t.Fatalf("Expected isEnabled = true, got <%v>", isEnabled)
	}

	if isActive != false {
		t.Fatalf("Expected isActive = false, got <%v>", isActive)
	}

	err = db.Close()
	if err != nil {
		t.Fatalf("sql.Close : Error : %s\n", err)
	}
}

func TestSelectNotNullAttribute(t *testing.T) {
	log.UseTestLogger(t)
	db, err := sql.Open("ramsql", "TestSelectNotNullAttribute")
	if err != nil {
		t.Fatalf("sql.Open : Error : %s\n", err)
	}
	defer db.Close()

	_, err = db.Exec(createTableQuery)
	if err != nil {
		t.Fatalf("Cannot create table: %s\n", err)
	}

	var id int
	var email sql.NullString
	var birthDate string

	// id=..., email=(NULL)SET, birth_date=(NOT-NULL)SET
	_, err = db.Exec("INSERT INTO user ('id', 'email', 'birth_date') VALUES (0, 'foo@bar.com', '2013-07-11T00:00:00Z')")
	if err != nil {
		t.Fatalf("Cannot insert into table: %s", err)
	}

	row := db.QueryRow("SELECT id, email, birth_date FROM user WHERE email = $1", "foo@bar.com")
	if row == nil {
		t.Fatalf("Cannot select row from table")
	}

	err = row.Scan(&id, &email, &birthDate)
	if err != nil {
		t.Fatalf("row.Scan: %s", err)
	}

	if id != 0 {
		t.Fatalf("Expected id = 0, got %d", id)
	}

	if !email.Valid || email.String != "foo@bar.com" {
		t.Fatalf("Expected email = 'foo@bar.com', got <%v>", email)
	}

	if birthDate != "2013-07-11T00:00:00Z" {
		t.Fatalf("Expected birth_date = '2013-07-11T00:00:00Z', got <%s>", birthDate)
	}

	// id=..., email=(NULL)UN-SET, birth_date=(NOT-NULL)SET
	_, err = db.Exec("INSERT INTO user ('id', 'birth_date') VALUES (1, '2013-07-12T00:00:00Z')")
	if err != nil {
		t.Fatalf("Cannot insert into table: %s", err)
	}

	row = db.QueryRow("SELECT id, email, birth_date FROM user WHERE id = $1", "1")
	if row == nil {
		t.Fatalf("Cannot select row from table")
	}

	err = row.Scan(&id, &email, &birthDate)
	if err != nil {
		t.Fatalf("row.Scan: %s", err)
	}

	if id != 1 {
		t.Fatalf("Expected id = 1, got %d", id)
	}

	if email.Valid {
		t.Fatalf("Expected email = nil, got <%v>", email)
	}

	if birthDate != "2013-07-12T00:00:00Z" {
		t.Fatalf("Expected birth_date = '2013-07-11T00:00:00Z', got <%s>", birthDate)
	}

	// id=..., email=(NULL)SET, birth_date=(NOT-NULL)UN-SET
	_, err = db.Exec("INSERT INTO user ('id', 'email') VALUES (2, 'bar@baz.com')")
	if err == nil {
		t.Fatalf("Should have failed to insert into table without specifying value for NOT NULL field 'email'")
	}
}

func TestCount(t *testing.T) {
	log.UseTestLogger(t)
	db, err := sql.Open("ramsql", "TestCount")
	if err != nil {
		t.Fatalf("sql.Open : Error : %s\n", err)
	}
	defer db.Close()

	batch := []string{
		`CREATE TABLE account (id BIGSERIAL, email TEXT)`,
		`INSERT INTO account (email) VALUES ("foo@bar.com")`,
		`INSERT INTO account (email) VALUES ("bar@bar.com")`,
		`INSERT INTO account (email) VALUES ("foobar@bar.com")`,
		`INSERT INTO account (email) VALUES ("babar@bar.com")`,
	}

	for _, b := range batch {
		_, err = db.Exec(b)
		if err != nil {
			t.Fatalf("sql.Exec: Error: %s\n", err)
		}
	}

	var count int64
	err = db.QueryRow(`SELECT COUNT(*) FROM account WHERE 1=1`).Scan(&count)
	if err != nil {
		t.Fatalf("cannot select COUNT of account: %s\n", err)
	}

	if count != 4 {
		t.Fatalf("Expected count to be 4, not %d", count)
	}

	err = db.QueryRow(`SELECT COUNT(i_dont_exist_lol) FROM account WHERE 1=1`).Scan(&count)
	if err == nil {
		t.Fatalf("Expected an error from a non existing attribute")
	}

}
