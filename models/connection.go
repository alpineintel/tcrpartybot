package models

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"os"
)

var pool *sqlx.DB = nil

func GetDBSession() *sqlx.DB {
	if pool != nil {
		return pool
	}

	session, err := sqlx.Open("postgres", os.Getenv("DATABASE"))

	if err != nil {
		log.Fatalf("Could not connect to specified database: %s", err)
	}

	pool = session
	return session
}
