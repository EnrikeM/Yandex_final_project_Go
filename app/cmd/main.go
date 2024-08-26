package main

import (
	"fmt"
	"log"

	_ "modernc.org/sqlite"

	"github.com/EnrikeM/Yandex_final_project_Go/app/internal/config"
	"github.com/EnrikeM/Yandex_final_project_Go/app/internal/httpsrv"
	"github.com/EnrikeM/Yandex_final_project_Go/app/internal/storage"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	config, err := config.New()
	if err != nil {
		return fmt.Errorf("error loading configuration: %v", err)
	}

	dbParams := storage.New(nil, *config)
	if err := dbParams.NewConnection(); err != nil {
		return fmt.Errorf("error connecting to database: %v", err)
	}

	api := httpsrv.NewAPI(dbParams.DB, *config)

	if err := api.Start(); err != nil {
		return fmt.Errorf("error starting API server: %v", err)
	}

	defer func() {
		if err := dbParams.DB.Close(); err != nil {
			log.Fatalf("error closing database: %v", err)
		}
	}()

	return nil
}
