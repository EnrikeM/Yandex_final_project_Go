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
	fmt.Println(dbParams.Config.TODO_DBFILE)

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
