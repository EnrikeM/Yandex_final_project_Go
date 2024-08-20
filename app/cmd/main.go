package main

import (
	"fmt"
	"log"
	"net/http"

	_ "modernc.org/sqlite"

	"github.com/EnrikeM/Yandex_final_project_Go/app/internal/config"
	"github.com/EnrikeM/Yandex_final_project_Go/app/internal/storage"
	"github.com/go-chi/chi"
)

var webDir = "web"

func main() {

	config := config.New()

	storage.New(config)

	r := chi.NewRouter()
	r.Handle("/*", http.StripPrefix("/", http.FileServer(http.Dir(webDir))))

	err := http.ListenAndServe(fmt.Sprintf(":%s", config.TODO_PORT), r)
	if err != nil {
		log.Fatal(fmt.Errorf("error starting server %w", err))
	}

}
