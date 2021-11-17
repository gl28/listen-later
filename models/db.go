package models

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var db *sql.DB

func Init() *sql.DB {
	connStr := "***REMOVED***"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	return db
}