package parser

import (
	"testing"

	"github.com/mlhoyt/ramsql/engine/log"
)

func TestParserCreateTableNoAttrConstraints(t *testing.T) {
	query := `CREATE TABLE account (id INT, email TEXT)`
	parse(query, 1, t)
}

func TestParserCreateTableWithPrimaryKeyConstraint(t *testing.T) {
	query := `CREATE TABLE account (id INT PRIMARY KEY, email TEXT)`
	parse(query, 1, t)
}

func TestParserCreateTableWithAutoIncrementConstraint(t *testing.T) {
	query := `CREATE TABLE account (id INT AUTOINCREMENT, email TEXT)`
	parse(query, 1, t)
}

func TestParserCreateTableWithOtherAutoIncrementConstraint(t *testing.T) {
	query := `CREATE TABLE account (id INT AUTO_INCREMENT, email TEXT)`
	parse(query, 1, t)
}

func TestParserCreateTableWithNotNullConstraint(t *testing.T) {
	query := `CREATE TABLE account (id INT PRIMARY KEY AUTO_INCREMENT, email TEXT NOT NULL)`
	parse(query, 1, t)
}

func TestParserCreateTableWithNullConstraint(t *testing.T) {
	query := `CREATE TABLE account (id INT PRIMARY KEY AUTO_INCREMENT, email TEXT NULL)`
	parse(query, 1, t)
}

func TestParserCreateTableNotNullPrimaryKeyConstraints(t *testing.T) {
	query := `CREATE TABLE account (id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY, email TEXT)`
	parse(query, 1, t)
}

func TestParserCreateTableWithDefaultOnUpdateConstraints(t *testing.T) {
	query := `CREATE TABLE account (id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY, modified_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP)`
	parse(query, 1, t)
}

func TestParserCreateTableWithConstraintPrimaryKey(t *testing.T) {
	query := `CREATE TABLE account (id BIGINT NOT NULL AUTO_INCREMENT, name TEXT NOT NULL, PRIMARY KEY (id))`
	parse(query, 1, t)
}

func TestParserCreateTableWithConstraintUniqueIndex(t *testing.T) {
	query := `CREATE TABLE account (id BIGINT NOT NULL AUTO_INCREMENT, name TEXT NOT NULL, UNIQUE INDEX (id))`
	parse(query, 1, t)
}

func TestParserCreateTableWithConstraintUniqueKey(t *testing.T) {
	query := `CREATE TABLE account (id BIGINT NOT NULL AUTO_INCREMENT, name TEXT NOT NULL, UNIQUE KEY (id))`
	parse(query, 1, t)
}

func TestParserCreateTableWithConstraintIndexNoNameNoType(t *testing.T) {
	query := `CREATE TABLE account (id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY, name TEXT NOT NULL, INDEX (name))`
	parse(query, 1, t)
}

func TestParserCreateTableWithConstraintIndexNoType(t *testing.T) {
	query := `CREATE TABLE account (id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY, name TEXT NOT NULL, INDEX name_idx (name))`
	parse(query, 1, t)
}

func TestParserCreateTableWithConstraintIndexNoNameTypeBtree(t *testing.T) {
	query := `CREATE TABLE account (id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY, name TEXT NOT NULL, INDEX USING BTREE (name))`
	parse(query, 1, t)
}

func TestParserCreateTableWithConstraintIndexNoNameTypeHash(t *testing.T) {
	query := `CREATE TABLE account (id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY, name TEXT NOT NULL, INDEX USING HASH (name))`
	parse(query, 1, t)
}

func TestParserCreateTableWithConstraintIndexAll(t *testing.T) {
	query := `CREATE TABLE account (id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY, name TEXT NOT NULL, INDEX name_idx USING BTREE (name))`
	parse(query, 1, t)
}

func TestParserCreateTableWithConstraintForeignKeyNoName(t *testing.T) {
	query := `CREATE TABLE account (id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY, name TEXT NOT NULL, FOREIGN KEY (name))`
	parse(query, 1, t)
}

func TestParserCreateTableWithConstraintForeignKey(t *testing.T) {
	query := `CREATE TABLE account (id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY, name TEXT NOT NULL, FOREIGN KEY name_idx (name))`
	parse(query, 1, t)
}

func TestParserCreateTableWithConstraintForeignKeyWithReferences(t *testing.T) {
	query := `CREATE TABLE account (id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY, name TEXT NOT NULL, FOREIGN KEY name_idx (name) REFERENCES user (name) ON DELETE CASCADE)`
	parse(query, 1, t)
}

func TestParserCreateTableWithShadowedKeyword(t *testing.T) {
	query := `CREATE TABLE policy (id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY, ` + "`action`" + ` VARCHAR(128) NOT NULL)`
	parse(query, 1, t)
}

func TestParserCreateTableWithBooleanField(t *testing.T) {
	query := `CREATE TABLE account (id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY, is_enabled BOOLEAN NOT NULL)`
	parse(query, 1, t)
}

func TestParserInsertTableWithBooleanFieldFalse(t *testing.T) {
	query := `INSERT INTO account (id, is_enabled) VALUES (0, false)`
	parse(query, 1, t)
}

func TestParserInsertTableWithBooleanFieldTrue(t *testing.T) {
	query := `INSERT INTO account (id, is_enabled) VALUES (0, true)`
	parse(query, 1, t)
}

func TestParserMultipleInstructions(t *testing.T) {
	query := `CREATE TABLE account (id INT, email TEXT);CREATE TABLE user (id INT, email TEXT)`
	parse(query, 2, t)
}

// func TestParserLowerCase(t *testing.T) {
// 	query := `create table account (id INT PRIMARY KEY NOT NULL)`
// 	parse(query, 1, t)
// }

func TestParserComplete(t *testing.T) {
	query := `CREATE TABLE user
	(
    	id INT PRIMARY KEY,
	    last_name TEXT,
	    first_name TEXT,
	    email TEXT,
	    birth_date DATE,
	    country TEXT,
	    town TEXT,
	    zip_code TEXT
	)`
	parse(query, 1, t)
}

func TestParserCompleteWithBacktickQuotes(t *testing.T) {
	query := `CREATE TABLE ` + "`" + `user` + "`" + `
	(
		` + "`" + `id` + "`" + ` INT PRIMARY KEY,
		` + "`" + `last_name` + "`" + ` TEXT,
		` + "`" + `first_name` + "`" + ` TEXT,
		` + "`" + `email` + "`" + ` TEXT,
		` + "`" + `birth_date` + "`" + ` DATE,
		` + "`" + `country` + "`" + ` TEXT,
		` + "`" + `town` + "`" + ` TEXT,
		` + "`" + `zip_code` + "`" + ` TEXT
	)`
	parse(query, 1, t)
}

// func TestParserCreateTableWithVarchar(t *testing.T) {
// 	query := `CREATE TABLE user
// 	(
//     	id INT PRIMARY KEY,
// 	    last_name VARCHAR(100)
// 	)`
// 	parse(query, 1, t)
// }

func TestParserCreateWithTableOptions1O1(t *testing.T) {
	query := `CREATE TABLE user ( id INT PRIMARY KEY, email TEXT, birth_date DATE ) ENGINE=InnoDB`
	parse(query, 1, t)
}

func TestParserCreateWithTableOptions1O2(t *testing.T) {
	query := `CREATE TABLE user ( id INT PRIMARY KEY, email TEXT, birth_date DATE ) ENGINE InnoDB`
	parse(query, 1, t)
}

func TestParserCreateWithTableOptions2aO1(t *testing.T) {
	query := `CREATE TABLE user ( id INT PRIMARY KEY, email TEXT, birth_date DATE ) DEFAULT CHARSET=utf8`
	parse(query, 1, t)
}

func TestParserCreateWithTableOptions2aO2(t *testing.T) {
	query := `CREATE TABLE user ( id INT PRIMARY KEY, email TEXT, birth_date DATE ) DEFAULT CHARSET utf8`
	parse(query, 1, t)
}

func TestParserCreateWithTableOptions2bO1(t *testing.T) {
	query := `CREATE TABLE user ( id INT PRIMARY KEY, email TEXT, birth_date DATE ) CHARSET=utf8`
	parse(query, 1, t)
}

func TestParserCreateWithTableOptions2bO2(t *testing.T) {
	query := `CREATE TABLE user ( id INT PRIMARY KEY, email TEXT, birth_date DATE ) CHARSET utf8`
	parse(query, 1, t)
}

func TestParserCreateWithTableOptions3aO1(t *testing.T) {
	query := `CREATE TABLE user ( id INT PRIMARY KEY, email TEXT, birth_date DATE ) DEFAULT CHARACTER SET=utf8`
	parse(query, 1, t)
}

func TestParserCreateWithTableOptions3aO2(t *testing.T) {
	query := `CREATE TABLE user ( id INT PRIMARY KEY, email TEXT, birth_date DATE ) DEFAULT CHARACTER SET utf8`
	parse(query, 1, t)
}

func TestParserCreateWithTableOptions3bO1(t *testing.T) {
	query := `CREATE TABLE user ( id INT PRIMARY KEY, email TEXT, birth_date DATE ) CHARACTER SET=utf8`
	parse(query, 1, t)
}

func TestParserCreateWithTableOptions3bO2(t *testing.T) {
	query := `CREATE TABLE user ( id INT PRIMARY KEY, email TEXT, birth_date DATE ) CHARACTER SET utf8`
	parse(query, 1, t)
}

func TestParserDropTableSimple(t *testing.T) {
	query := `DROP TABLE myTableName`
	parse(query, 1, t)
}

func TestParserDropTableIfExists(t *testing.T) {
	query := `DROP TABLE IF EXISTS myTableName`
	parse(query, 1, t)
}

func TestSelectStar(t *testing.T) {
	query := `SELECT * FROM account WHERE email = 'foo@bar.com'`
	parse(query, 1, t)
}

func TestSelectMultipleAttribute(t *testing.T) {
	query := `SELECT id, email FROM account WHERE email = 'foo@bar.com'`
	parse(query, 1, t)
}

func TestSelectOneAttribute(t *testing.T) {
	query := `SELECT id FROM account WHERE email = 'foo@bar.com'`
	parse(query, 1, t)
}

func TestSelectAttributeWithTable(t *testing.T) {
	query := `SELECT account.id FROM account WHERE email = 'foo@bar.com'`
	parse(query, 1, t)
}

func TestSelectAttributeWithQuotedTable(t *testing.T) {
	query := `SELECT "account".id FROM account WHERE email = 'foo@bar.com'`
	parse(query, 1, t)
}

func TestSelectAttributeWithBacktickQuotedTable(t *testing.T) {
	query := `SELECT ` + "`" + `account` + "`" + `.id FROM account WHERE email = 'foo@bar.com'`
	parse(query, 1, t)
}

func TestSelectAllFromTable(t *testing.T) {
	query := `SELECT "account".* FROM account WHERE email = 'foo@bar.com'`
	parse(query, 1, t)
}

func TestSelectOnePredicate(t *testing.T) {
	query := `SELECT * FROM account WHERE 1`
	parse(query, 1, t)
}

func TestSelectQuotedTableName(t *testing.T) {
	query := `SELECT * FROM "account" WHERE 1`
	parse(query, 1, t)

	query = `SELECT * FROM "account"`
	parse(query, 1, t)
}

func TestSelectBacktickQuotedTableName(t *testing.T) {
	query := `SELECT * FROM ` + "`" + `account` + "`" + ` WHERE 1`
	parse(query, 1, t)

	query = `SELECT * FROM ` + "`" + `account` + "`" + ``
	parse(query, 1, t)
}

func TestSelectJoin(t *testing.T) {
	query := `SELECT address.* FROM address
	JOIN user_addresses ON address.id=user_addresses.address_id
	WHERE user_addresses.user_id=1`
	parse(query, 1, t)
}

func TestInsertMinimal(t *testing.T) {
	query := `INSERT INTO account ('email', 'password', 'age') VALUES ('foo@bar.com', 'tititoto', '4')`
	parse(query, 1, t)
}

func TestInsertNumber(t *testing.T) {
	query := `INSERT INTO account ('email', 'password', 'age') VALUES ('foo@bar.com', 'tititoto', 4)`
	parse(query, 1, t)
}

func TestInsertNumberWithQuote(t *testing.T) {
	query := `INSERT INTO "account" ('email', 'password', 'age') VALUES ('foo@bar.com', 'tititoto', 4)`
	parse(query, 1, t)
}

func TestInsertNumberWithBacktickQuote(t *testing.T) {
	query := `INSERT INTO ` + "`" + `account` + "`" + ` ('email', 'password', 'age') VALUES ('foo@bar.com', 'tititoto', 4)`
	parse(query, 1, t)
}

func TestCreateTableWithKeywordName(t *testing.T) {
	query := `CREATE TABLE test ("id" bigserial not null primary key, "name" text, "key" text)`
	parse(query, 1, t)
}

// func TestInsertStringWithDoubleQuote(t *testing.T) {
// 	query := `insert into "posts" ("post_id","Created","Title","Body") values (null,12321123,"Hello world !","!");`
// 	parse(query, 1, t)
// }

func TestInsertStringWithSimpleQuote(t *testing.T) {
	query := `insert into "posts" ("post_id","Created","Title","Body") values (null,12321123,'Hello world !','!');`
	parse(query, 1, t)
}

// func TestInsertImplicitAttributes(t *testing.T) {
// 	query := `INSERT INTO account VALUES ('foo@bar.com', 'tititoto', 4)`
// 	parse(query, 1, t)
// }

func TestParseDelete(t *testing.T) {
	query := `delete from "posts"`
	parse(query, 1, t)
}

func TestParseUpdate(t *testing.T) {
	query := `UPDATE account SET email = 'roger@gmail.com' WHERE id = 2`
	parse(query, 1, t)
}

func TestUpdateMultipleAttributes(t *testing.T) {
	query := `update "posts" set "Created"=1435760856063203203, "Title"='Go 1.2 is better than ever', "Body"='Lorem ipsum lorem ipsum' where "post_id"=2`
	parse(query, 1, t)
}

func TestParseMultipleJoin(t *testing.T) {
	query := `SELECT group.id, user.username FROM group JOIN group_user ON group_user.group_id = group.id JOIN user ON user.id = group_user.user_id WHERE group.name = 1`
	parse(query, 1, t)
}

func TestParseMultipleOrderBy(t *testing.T) {
	query := `SELECT group.id, user.username FROM group JOIN group_user ON group_user.group_id = group.id JOIN user ON user.id = group_user.user_id WHERE group.name = 1 ORDER BY group.name, user.username ASC`
	parse(query, 1, t)
}

func TestSelectForUpdate(t *testing.T) {
	query := `SELECT * FROM user WHERE user.id = 1 FOR UPDATE`

	parse(query, 1, t)
}

func TestCreateDefault(t *testing.T) {
	query := `CREATE TABLE foo (bar BIGINT, riri TEXT, fifi BOOLEAN NOT NULL DEFAULT false)`

	parse(query, 1, t)
}

func TestCreateDefaultNumerical(t *testing.T) {
	query := `CREATE TABLE foo (bar BIGINT, riri TEXT, fifi BIGINT NOT NULL DEFAULT 0)`

	parse(query, 1, t)
}

func TestCreateWithTimestamp(t *testing.T) {
	query := `CREATE TABLE IF NOT EXISTS "pokemon" (id BIGSERIAL PRIMARY KEY, name TEXT, type TEXT, seen TIMESTAMP WITH TIME ZONE)`

	parse(query, 1, t)
}

func TestCreateDefaultTimestamp(t *testing.T) {
	query := `CREATE TABLE IF NOT EXISTS "pokemon" (id BIGSERIAL PRIMARY KEY, name TEXT, type TEXT, seen TIMESTAMP WITH TIME ZONE DEFAULT LOCALTIMESTAMP)`

	parse(query, 1, t)
}

func TestCreateNumberInNames(t *testing.T) {
	query := `CREATE TABLE IF NOT EXISTS "pokemon" (id BIGSERIAL PRIMARY KEY, name TEXT, type TEXT, md5sum TEXT)`

	parse(query, 1, t)
}

func TestOffset(t *testing.T) {
	query := `SELECT * FROM mytable LIMIT 1 OFFSET 0`

	parse(query, 1, t)
}

func TestUnique(t *testing.T) {
	queries := []string{
		`CREATE TABLE pokemon (id BIGSERIAL, name TEXT UNIQUE NOT NULL)`,
		`CREATE TABLE pokemon (id BIGSERIAL, name TEXT NOT NULL UNIQUE)`,
		`CREATE TABLE pokemon_name (id BIGINT, name VARCHAR(255) PRIMARY KEY NOT NULL UNIQUE)`,
	}

	for _, q := range queries {
		parse(q, 1, t)
	}
}

func parse(query string, instructionCount int, t *testing.T) []Instruction {
	log.UseTestLogger(t)

	lexer := lexer{}
	tokens, err := lexer.lex([]byte(query))
	if err != nil {
		t.Fatalf("Cannot lex <%s> string: %s", query, err)
	}

	parser := NewParser(tokens)
	instructions, err := parser.parse()
	if err != nil {
		t.Fatalf("Cannot parse tokens from '%s': %s", query, err)
	}

	if len(instructions) != instructionCount {
		t.Fatalf("Should have parsed %d instructions, got %d", instructionCount, len(instructions))
	}

	return instructions
}
