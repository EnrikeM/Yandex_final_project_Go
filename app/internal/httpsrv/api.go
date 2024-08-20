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

func NewAPI(config config.Config) *API {
	return &API{
		r:      chi.NewRouter(),
		config: config,
	}
}

func (a *API) Register(r chi.Router) {
	r.Handle("/*", http.FileServer(http.Dir(a.config.WEB_DIR)))
	// r.Get("/api/nextdate", a.GetHandler)
}

func (a *API) Start() {
	a.Register(a.r)
	fmt.Println("server start")

	err := http.ListenAndServe(fmt.Sprintf(":%s", a.config.TODO_PORT), a.r)
	if err != nil {
		log.Fatal(fmt.Errorf("error starting server %w", err))
	}
}
