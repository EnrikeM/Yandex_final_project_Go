package storage

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/EnrikeM/Yandex_final_project_Go/app/internal/config"
)

type DBParams struct {
	DB     *sql.DB
	Config config.Config
}

func New(db *sql.DB, config config.Config) DBParams {
	return DBParams{
		DB:     db,
		Config: config,
	}
}

func (dbParams *DBParams) NewConnection() error {
	dbFile := dbParams.Config.TODO_DBFILE

	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return dbParams.createDatabase(dbFile)
	}

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("error opening db: %w", err)
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("error pinging db: %w", err)
	}

	dbParams.DB = db
	return nil
}

func (dbParams *DBParams) createDatabase(dbFile string) error {
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("error opening db: %w", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date   VARCHAR(32) NOT NULL DEFAULT "0001-01-01",
		title   VARCHAR(64) NOT NULL DEFAULT "",
		comment VARCHAR(128) NOT NULL DEFAULT "",
		repeat  VARCHAR(128) NOT NULL DEFAULT ""
	);
	CREATE INDEX scheduler_date ON scheduler (date);`)
	if err != nil {
		return fmt.Errorf("error creating db: %w", err)
	}

	dbParams.DB = db
	return nil
}
