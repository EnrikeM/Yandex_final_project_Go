package httpsrv

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/EnrikeM/Yandex_final_project_Go/app/internal/config"
	"github.com/go-chi/chi"
)

type API struct {
	DB     *sql.DB
	config config.Config
	r      chi.Router
}

func NewAPI(db *sql.DB, config config.Config) *API {
	return &API{
		r:      chi.NewRouter(),
		config: config,
		DB:     db,
	}
}

func (a *API) Register(r chi.Router) {

	r.Get("/api/nextdate", a.GetNextDateHandler)
	r.Post("/api/task", a.PostTaskHandler)
	r.Get("/api/tasks", a.GetTaskHandler)

	r.Handle("/*", http.FileServer(http.Dir(a.config.WEB_DIR)))
}

func (a *API) Start() error {
	a.Register(a.r)
	fmt.Println("server start")

	err := http.ListenAndServe(fmt.Sprintf(":%s", a.config.TODO_PORT), a.r)
	if err != nil {
		log.Fatal(fmt.Errorf("error starting server %w", err))
	}
	return nil
}
