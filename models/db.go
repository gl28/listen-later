package models

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

func Init() (*sql.DB, error) {
	connStr := os.Getenv("LISTEN_LATER_DB")
	var err error
	db, err = sql.Open("postgres", connStr)
	
	return db, err
}