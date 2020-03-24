package main

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/pkg/errors"
	"time"
)

var db *sql.DB
var st *sql.Stmt

func PrintDBStats() {
	stats := db.Stats()
	fmt.Printf("OpenConnections: %v, Idle: %v, InUse: %v, IdleClosed: %v, LifetimeClosed: %v\n", stats.OpenConnections, stats.Idle, stats.InUse, stats.MaxIdleClosed, stats.MaxLifetimeClosed)
}

func Execute() {
	var pid string
	if err := st.QueryRow().Scan(&pid); err != nil {
		panic(err)
	}

	fmt.Println("PID: ", pid)

	PrintDBStats()
}

func main() {
	var err error
	db, err = sql.Open("pgx", "postgres://postgres@localhost:5432/learn?sslmode=disable")
	if err != nil {
		panic(errors.Wrapf(err, "Error opening database"))
	}

	db.SetMaxIdleConns(2)
	db.SetConnMaxLifetime(time.Second * 300)

	st, err = db.Prepare("SELECT pg_backend_pid()")
	if err != nil {
		panic(errors.Wrapf(err, "Error preparing statement"))
	}

	for i := 0; i < 3; i++ {
		Execute()
	}

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	Execute()

	if err := tx.Rollback(); err != nil {
		panic(err)
	}

	time.Sleep(time.Second * 10)
	PrintDBStats()
	
	Execute()
}
