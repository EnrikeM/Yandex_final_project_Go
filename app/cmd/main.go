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
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	config, err := config.New()
	if err != nil {
		return fmt.Errorf("error loading configuration: %v", err)
	}

	storage := storage.New(nil, config)
	if err := storage.NewConnection(); err != nil {
		return fmt.Errorf("error connecting to database: %v", err)
	}

	api := httpsrv.NewAPI(config, storage)

	if err := api.Start(); err != nil {
		return fmt.Errorf("error starting API server: %v", err)
	}

	return nil
}
