package httpsrv

import (
	"fmt"
	"log"
	"net/http"

	"github.com/EnrikeM/Yandex_final_project_Go/app/internal/config"
	"github.com/go-chi/chi"
)

type API struct {
	config config.Config
	r      chi.Router
}

func New(r chi.Router, config config.Config) *API {
	return &API{
		r:      r,
		config: config,
	}
}

func (a *API) Register(r chi.Router) {
	r.Handle("/*", http.StripPrefix("/", http.FileServer(http.Dir(a.config.WEB_DIR))))
	// r.Get("/api/nextdate", a.GetHandler)
}

func (a *API) Start() {
	err := http.ListenAndServe(fmt.Sprintf(":%s", a.config.TODO_PORT), a.r)
	if err != nil {
		log.Fatal(fmt.Errorf("error starting server %w", err))
	}
}
