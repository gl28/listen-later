package models

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var db *sql.DB

func Init() *sql.DB {
	connStr := "postgres://gdikpsow:POMc3kLSykI29JU87h5nzUbvVWAZMqKr@kashin.db.elephantsql.com/gdikpsow"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	return db
}