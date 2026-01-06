package database

import (
	"database/sql"
	"log"

	"github.com/yogusita/to-adhdo/env"
	_ "modernc.org/sqlite"
)

const DEFAULT_DB_FILE_NAME = "main.db"

var db_client *sql.DB = nil

func GetDatabase() (*sql.DB, error) {
	if db_client != nil || (db_client != nil && db_client.Stats().OpenConnections == 0) {
		log.Printf("returning existing db_client: %v\n", db_client)
		return db_client, nil
	}

	data_souce_name, _ := env.LookupEnvWithDefault("DB_NAME", DEFAULT_DB_FILE_NAME)

	log.Printf("opening new client for \"%s\"", data_souce_name)

	var err error
	db_client, err = sql.Open("sqlite", data_souce_name)

	log.Println("client opened")

	if err != nil {
		log.Printf("client open error: %s\n", err.Error())
		return nil, err
	}

	log.Println("pinging")

	err = db_client.Ping()

	if err != nil {
		log.Printf("ping error: %s\n", err.Error())
		return nil, err
	}

	log.Printf("returning new db_client: %v", db_client)

	return db_client, nil
}
