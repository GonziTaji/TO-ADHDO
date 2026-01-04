package db

import (
	"database/sql"

	"github.com/yogusita/to-adhdo/env"
)

const DEFAULT_DB_FILE_NAME = "database"

var open_database *sql.DB = nil

func GetDatabase() (*sql.DB, error) {
	if open_database != nil {
		return open_database, nil
	}

	data_souce_name, _ := env.LookupEnvWithDefault("DB_NAME", DEFAULT_DB_FILE_NAME)

	var err error
	open_database, err = sql.Open("sqlite", data_souce_name)

	if err != nil {
		return nil, err
	}

	return open_database, nil
}
