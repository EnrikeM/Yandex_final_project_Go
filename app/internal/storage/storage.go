package storage

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/EnrikeM/Yandex_final_project_Go/app/internal/config"
)

func New(config *config.Config) {

	dbFile := config.TODO_DBFILE

	_, err := os.Stat(dbFile)

	var install bool
	if err != nil {
		if os.IsNotExist(err) {
			install = true
		} else {
			log.Fatal(fmt.Errorf("error checking file: %w", err))
		}

	}

	if install {
		db, err := sql.Open("sqlite", dbFile)
		if err != nil {
			log.Fatal(fmt.Errorf("error opening db: %w", err))
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
			log.Fatal(fmt.Errorf("error creating db: %w", err))
		}

	}
}
