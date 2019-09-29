package main

import (
	"database/sql"
	"fmt"

	"github.com/mlhoyt/ramsql/cli"
	_ "github.com/mlhoyt/ramsql/driver"
)

func main() {
	db, err := sql.Open("ramsql", "")
	if err != nil {
		fmt.Printf("Error : cannot open connection : %s\n", err)
		return
	}
	cli.Run(db)
}
