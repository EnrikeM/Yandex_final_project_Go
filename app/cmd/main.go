package main

import (
	_ "modernc.org/sqlite"

	"github.com/EnrikeM/Yandex_final_project_Go/app/internal/config"
	"github.com/EnrikeM/Yandex_final_project_Go/app/internal/httpsrv"
	"github.com/EnrikeM/Yandex_final_project_Go/app/internal/storage"
	"github.com/go-chi/chi"
)

func main() {

	config := config.New()
	storage.New(config)
	r := chi.NewRouter()
	api := httpsrv.New(r, *config)
	api.Register(r)
	api.Start()

}
