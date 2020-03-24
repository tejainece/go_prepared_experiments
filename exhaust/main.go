package main

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/pkg/errors"
)

var db *sql.DB

func main() {
	var err error
	
	db, err = sql.Open("pgx", "postgres://postgres@localhost:5432/learn?sslmode=disable")
	if err != nil {
		panic(errors.Wrapf(err, "error opening database"))
	}
	
	count := 100000
	
	sts  := make([]*sql.Stmt, count)
	
	for i := 0; i < count; i++ {
		st, err := db.Prepare(fmt.Sprintf("SELECT pg_backend_pid(), %v", i))
		if err != nil {
			panic(errors.Wrapf(err, "error preparing statement"))
		}
		
		var pid, data int
		if err := st.QueryRow().Scan(&pid, &data); err != nil {
			panic(err)
		}
		
		fmt.Println(pid, data)

		sts[i] = st
	}
	
	for i := 0; i < count; i++ {
		st := sts[i]

		var pid, data int
		if err := st.QueryRow().Scan(&pid, &data); err != nil {
			panic(err)
		}

		fmt.Println(pid, data)

		sts[i] = st
	}
}
