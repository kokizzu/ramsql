# Design / Implementation Notes

## Call Stack

```
+-- engine.Engine::handleConnection(conn protocol.EngineConn) none
    +-- protocol.EngineConn::ReadStatement() -> (stmt string)
    +-- parser.ParseInstruction(stmt string) -> (instructions []parser.Instruction, err error)
    +-- engine.Engine::executeQueries(instructions []parser.Instruction, conn protocol.EngineConn) -> (err error)
        +-- engine.Engine::executeQuery(instruction parser.Instruction, conn protocol.EngineConn) -> (err error)
            +-- engine.executor(*engine.Engine, *parser.Decl, protocol.EngineConn) -> (err error)
                * NOTE: There are engine.{...}Executor functions that do NOT meet the executor type (joinExecutor, orderbyExecutor, inExecutor, whereExecutor, fromExecutor)
```

## Parser Instruction

### parser.Instruction

TODO: Need to document the structure of parser.Instruction type

### parser.Decl

An AST structure for representing a SQL statement
```go
type Decl struct {
  Token  int
  Lexeme string
  Decl   []*Decl
}
```

## Issues

### DOES work

```sql
CREATE TABLE IF NOT EXISTS customer (
  customer_id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(128) NOT NULL,
  customer_key VARCHAR(16) NOT NULL,
  encryption_key VARCHAR(32) NOT NULL,
  integration_key VARCHAR(36) NOT NULL,
  first_name VARCHAR(64) NOT NULL,
  last_name VARCHAR(64) NOT NULL,
  street_line1 VARCHAR(128) NULL,
  street_line2 VARCHAR(128) NULL,
  city VARCHAR(128) NULL,
  state VARCHAR(32) NULL,
  postal_code VARCHAR(16) NULL,
  work_phone VARCHAR(32) NULL,
  mobile_phone VARCHAR(32) NULL,
  email VARCHAR(128),
  image_url VARCHAR(512),
  created_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  modified_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  created_by VARCHAR(64) NULL,
  modified_by VARCHAR(64) NULL,
  INDEX idx_customer_customer_key (customer_key)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
```

### DOES work

```sql
CREATE TABLE IF NOT EXISTS user (
  user_id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  customer_id BIGINT NOT NULL,
  first_name VARCHAR(64),
  last_name VARCHAR(64),
  email VARCHAR(256) NOT NULL,
  password VARCHAR(1024) NOT NULL,
  is_active BOOLEAN NOT NULL,
  is_enabled BOOLEAN NOT NULL,
  is_first_time tinyint(1) NOT NULL DEFAULT 1,
  image_url VARCHAR(512),
  created_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  modified_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  created_by VARCHAR(64) NULL,
  modified_by VARCHAR(64) NULL,
  INDEX idx_user_customer_id (customer_id),
  INDEX idx_user_email (email),
  FOREIGN KEY (customer_id)
    REFERENCES customer(customer_id)
    ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
```

### DOES work

```sql
CREATE TABLE IF NOT EXISTS standard_alpha (
  standard_id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(256) NOT NULL,
  version VARCHAR(16) NOT NULL,
  month INT NOT NULL,
  year INT NOT NULL,
  created_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  modified_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  created_by VARCHAR(64) NOT NULL,
  modified_by VARCHAR(64) NOT NULL,
  modified_by_id BIGINT(20) NOT NULL,
  FOREIGN KEY (modified_by_id)
    REFERENCES user (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
```

### does NOT work

```sql
CREATE TABLE IF NOT EXISTS section_alpha (
  section_id BIGINT(20) NOT NULL AUTO_INCREMENT PRIMARY KEY,
  standard_id BIGINT(20) NOT NULL,
  title VARCHAR(256) NOT NULL,
  text TEXT NULL DEFAULT NULL,
  number VARCHAR(32) NOT NULL,
  abbreviation VARCHAR(8) NOT NULL,
  color VARCHAR(8) NOT NULL,
  created_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  modified_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  created_by VARCHAR(64) NOT NULL,
  modified_by VARCHAR(64) NOT NULL,
  modified_by_id BIGINT(20) NOT NULL,
  INDEX idx_section_standard_id (standard_id ASC),
  FOREIGN KEY (standard_id)
    REFERENCES standard_alpha (standard_id)
    ON DELETE CASCADE,
  FOREIGN KEY (modified_by_id)
    REFERENCES user (user_id)
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4
```

It appears this is the line that fails to parse:

```sql
  text TEXT NULL DEFAULT NULL,
```
