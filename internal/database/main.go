package server

import (
	"database/sql"
	"log"

	"github.com/yogusita/to-adhdo/env"
)

func newDB() (*sql.DB, error) {
	data_souce_name, _ := env.LookupEnvWithDefault("DB_NAME", "database/main.db")

	db, err := sql.Open("sqlite", data_souce_name)

	if err != nil {
		log.Printf("client open error: %s\n", err.Error())
		return nil, err
	}

	if err := db.Ping(); err != nil {
		log.Printf("server ping error: %s\n", err.Error())
		return nil, err
	}

	_, err = db.Exec("PRAGMA foreign_keys = ON;")

	if err != nil {
		log.Printf("pragma set error: %s\n", err.Error())
		return nil, err
	}

	return db, nil
}
