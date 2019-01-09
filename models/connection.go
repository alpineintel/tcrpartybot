package models

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"os"
)

var pool *sqlx.DB = nil

// GetDBSession returns the current active database connection pool or creates
// it if it doesn't already exist
func GetDBSession() *sqlx.DB {
	if pool != nil {
		return pool
	}

	session, err := sqlx.Open("postgres", os.Getenv("DATABASE_URL"))

	if err != nil {
		log.Fatalf("Could not connect to specified database: %s", err)
	}

	pool = session
	return session
}
