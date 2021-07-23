package engine_test

import (
	"database/sql"
	"testing"

	"github.com/kokizzu/ramsql/engine/log"

	_ "github.com/kokizzu/ramsql/driver"
)

func TestJoinOrderBy(t *testing.T) {
	log.UseTestLogger(t)

	db, err := sql.Open("ramsql", "TestJoinOrderBy")
	if err != nil {
		t.Fatalf("sql.Open: %s", err)
	}
	defer db.Close()

	init := []string{
		`CREATE TABLE user (id BIGSERIAL, name TEXT)`,
		`CREATE TABLE address (id BIGSERIAL, user_id INT, value TEXT)`,
		`INSERT INTO user (name) VALUES ($$riri$$)`,
		`INSERT INTO user (name) VALUES ($$fifi$$)`,
		`INSERT INTO user (name) VALUES ($$loulou$$)`,
		`INSERT INTO address (user_id, value) VALUES (1, 'rue du puit')`,
		`INSERT INTO address (user_id, value) VALUES (1, 'rue du désert')`,
		`INSERT INTO address (user_id, value) VALUES (3, 'rue du chemin')`,
		`INSERT INTO address (user_id, value) VALUES (2, 'boulevard du con')`,
	}
	for _, q := range init {
		_, err := db.Exec(q)
		if err != nil {
			t.Fatalf("Cannot initialize test: %s", err)
		}
	}

	query := `SELECT user.name, address.value 
			FROM user 
			JOIN address ON address.user_id = user.id
			WHERE user.id = $1
			ORDER BY address.value ASC`
	rows, err := db.Query(query, 1)
	if err != nil {
		t.Fatalf("Cannot select with joined order by: %s", err)
	}
	defer rows.Close()

	n := 0
	for rows.Next() {
		n++
	}
	if n != 2 {
		t.Fatalf("Expected 2 rows, got %d", n)
	}
}

func TestJoinAs(t *testing.T) {
	log.UseTestLogger(t)

	db, err := sql.Open("ramsql", "TestJoinAs")
	if err != nil {
		t.Fatalf("sql.Open: %s", err)
	}
	defer db.Close()

	init := []string{
		`CREATE TABLE user (id BIGINT PRIMARY KEY AUTO_INCREMENT, name TEXT)`,
		`CREATE TABLE address (id BIGINT PRIMARY KEY AUTO_INCREMENT, user_id INT, value TEXT)`,
		`INSERT INTO user (name) VALUES ('foo')`,
		`INSERT INTO user (name) VALUES ('bar')`,
		`INSERT INTO address (user_id, value) VALUES (1, '123 First St')`,
		`INSERT INTO address (user_id, value) VALUES (2, '456 Second St')`,
		`INSERT INTO address (user_id, value) VALUES (2, 'Apt Alpha')`,
	}
	for _, q := range init {
		_, err := db.Exec(q)
		if err != nil {
			t.Fatalf("Cannot initialize test: %s", err)
		}
	}

	query := `SELECT user.id AS user_id, user.name AS user_name, address.value AS user_address
			FROM user 
			JOIN address ON user.id = address.user_id
			WHERE user.id = ?`
	rows, err := db.Query(query, 2)
	if err != nil {
		t.Fatalf("Cannot select with joined as: %s", err)
	}
	defer rows.Close()

	n := 0
	for rows.Next() {
		var user_id int64
		var user_name string
		var user_address string

		err := rows.Scan(&user_id, &user_name, &user_address)
		if err != nil {
			t.Fatalf("Failed to scan row %d: %s", n, err)
		}

		n++
	}
	if n != 2 {
		t.Fatalf("Expected 2 rows, got %d", n)
	}
}

func TestJoinWithMixedAs(t *testing.T) {
	log.UseTestLogger(t)

	db, err := sql.Open("ramsql", "TestJoinWithMixedAs")
	if err != nil {
		t.Fatalf("sql.Open: %s", err)
	}
	defer db.Close()

	init := []string{
		`CREATE TABLE address (
       id BIGINT PRIMARY KEY AUTO_INCREMENT,
       value TEXT
     )`,
		`CREATE TABLE user (
       id BIGINT PRIMARY KEY AUTO_INCREMENT,
       name TEXT,
       address_id INT
     )`,
	}
	for _, q := range init {
		_, err := db.Exec(q)
		if err != nil {
			t.Fatalf("Cannot initialize test: %s", err)
		}
	}

	query := `SELECT
      user.id,
  		user.name,
  		address.id AS address_id,
  		address.value AS address_value
		FROM user 
		INNER JOIN address
		ON user.address_id = address.id`
	// WHERE user.id = ?`
	// rows, err := db.Query(query, 2)
	rows, err := db.Query(query)
	if err != nil {
		t.Fatalf("Cannot select with joined as: %s", err)
	}
	defer rows.Close()

	rowColumns, err := rows.Columns()
	if err != nil {
		t.Fatal("failed to get rows.Columns() for comparison")
	}
	expectedRowColumns := []string{"id", "name", "address_id", "address_value"}
	if len(rowColumns) != len(expectedRowColumns) {
		t.Fatalf("expected %d row columns, got %d", len(expectedRowColumns), len(rowColumns))
	}
	for i := range expectedRowColumns {
		if expectedRowColumns[i] != rowColumns[i] {
			t.Errorf("at position %d expected column name %q, got %q", i, expectedRowColumns[i], rowColumns[i])
		}
	}
}

func TestMultipleJoin(t *testing.T) {
	log.UseTestLogger(t)

	db, err := sql.Open("ramsql", "TestMultipleJoin")
	if err != nil {
		t.Fatalf("sql.Open: %s", err)
	}
	defer db.Close()

	init := []string{
		`CREATE TABLE user (id BIGSERIAL, name TEXT)`,
		`CREATE TABLE address (id BIGSERIAL, user_id INT, value TEXT)`,
		`CREATE TABLE user_group (id BIGSERIAL, user_id INT, name TEXT)`,
		`INSERT INTO user (name) VALUES ($$riri$$)`,
		`INSERT INTO user (name) VALUES ($$fifi$$)`,
		`INSERT INTO user (name) VALUES ($$loulou$$)`,
		`INSERT INTO address (user_id, value) VALUES (1, 'rue du puit')`,
		`INSERT INTO address (user_id, value) VALUES (1, 'rue du désert')`,
		`INSERT INTO address (user_id, value) VALUES (3, 'rue du chemin')`,
		`INSERT INTO address (user_id, value) VALUES (2, 'boulevard du con')`,
		`INSERT INTO user_group (user_id, name) VALUES (1, 'toto')`,
		`INSERT INTO user_group (user_id, name) VALUES (2, 'toto')`,
		`INSERT INTO user_group (user_id, name) VALUES (3, 'lonely')`,
		`INSERT INTO user_group (user_id, name) VALUES (1, 'cowboy')`,
	}
	for _, q := range init {
		_, err := db.Exec(q)
		if err != nil {
			t.Fatalf("Cannot initialize test: %s", err)
		}
	}

	query := `SELECT user.name, address.value
			FROM user
			JOIN address ON address.user_id = user.id
			JOIN user_group ON user_group.user_id = address.user_id
			WHERE user_group.name = $1
			ORDER BY address.value ASC`
	rows, err := db.Query(query, "toto")
	if err != nil {
		t.Fatalf("Cannot select with 3 tables joined: %s", err)
	}
	defer rows.Close()

	n := 0
	for rows.Next() {
		n++
	}
	if n != 3 {
		t.Fatalf("Expected 3 rows, got %d", n)
	}

}

func TestJoinGroup(t *testing.T) {
	log.UseTestLogger(t)

	db, err := sql.Open("ramsql", "TestJoinGroup")
	if err != nil {
		t.Fatalf("sql.Open: %s", err)
	}
	defer db.Close()

	init := []string{
		`CREATE TABLE user (id BIGSERIAL, name TEXT)`,
		`CREATE TABLE group (id BIGSERIAL, name TEXT)`,
		`CREATE TABLE user_group (id BIGSERIAL, user_id INT, group_id INT)`,
		`INSERT INTO user (name) VALUES ($$riri$$)`,
		`INSERT INTO user (name) VALUES ($$fifi$$)`,
		`INSERT INTO user (name) VALUES ($$loulou$$)`,
		`INSERT INTO group (name) VALUES ('cowboys')`,
		`INSERT INTO group (name) VALUES ('troopers')`,
		`INSERT INTO group (name) VALUES ('toys')`,
		`INSERT INTO user_group (user_id, group_id) VALUES (1, 1)`,
		`INSERT INTO user_group (user_id, group_id) VALUES (2, 1)`,
	}
	for _, q := range init {
		_, err := db.Exec(q)
		if err != nil {
			t.Fatalf("Cannot initialize test: %s", err)
		}
	}

	query := `SELECT user.name FROM user
			JOIN user_group ON user.id = user_group.user_id
			WHERE user_group.group_id = $1
			ORDER BY user.name ASC`
	rows, err := db.Query(query, 1)
	if err != nil {
		t.Fatalf("Cannot select joined: %s", err)
	}
	defer rows.Close()

	n := 0
	for rows.Next() {
		n++
	}
	if n != 2 {
		t.Fatalf("Expected 2 rows, got %d", n)
	}

}
