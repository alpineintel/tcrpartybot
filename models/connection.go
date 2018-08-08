package models

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

var session *sqlx.DB = nil

func GetDBSession() *sqlx.DB {
	if session != nil {
		return session
	}

	log.Println(os.Getenv("DATABASE_URI"))
	session, err := sqlx.Connect("sqlite3", os.Getenv("DATABASE_URI"))

	if err != nil {
		log.Fatal("Could not connect to specified database")
	}

	return session
}
