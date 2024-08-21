package main

import (
	_ "modernc.org/sqlite"

	"github.com/EnrikeM/Yandex_final_project_Go/app/internal/config"
	"github.com/EnrikeM/Yandex_final_project_Go/app/internal/httpsrv"
	"github.com/EnrikeM/Yandex_final_project_Go/app/internal/storage"
)

func main() {

	config := config.New()
	storage.New(config)
	api := httpsrv.NewAPI(*config)

	api.Start()
}
