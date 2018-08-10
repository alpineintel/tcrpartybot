package models

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"time"
)

var session *sqlx.DB = nil

type NullTime struct {
	Time  time.Time
	Valid bool
}

func GetDBSession() *sqlx.DB {
	if session != nil {
		return session
	}

	session, err := sqlx.Connect("sqlite3", os.Getenv("DATABASE_URI"))

	if err != nil {
		log.Fatal("Could not connect to specified database")
	}

	return session
}