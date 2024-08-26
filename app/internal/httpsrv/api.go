package httpsrv

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/EnrikeM/Yandex_final_project_Go/app/internal/config"
	"github.com/go-chi/chi"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/EnrikeM/Yandex_final_project_Go/docs"
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

	r.Post("/api/signin", http.HandlerFunc(a.auth(http.HandlerFunc(a.signInHandler)).ServeHTTP))
	r.Route("/api/task", func(r chi.Router) {
		r.Use(a.auth)
		r.Post("/", a.postTaskHandler)
		r.Get("/", a.getTaskHandler)
		r.Put("/", a.updateTaskHandler)
		r.Post("/done", a.postDoneHandler)
		r.Delete("/", a.deleteTaskHandler)
	})

	r.Route("/api/tasks", func(r chi.Router) {
		r.Get("/", a.getTasksHandler)
	})

	r.Get("/api/nextdate", a.getNextDateHandler)

	r.Get("/swagger/*", httpSwagger.WrapHandler)
	r.Handle("/*", http.FileServer(http.Dir(a.config.WEB_DIR)))
}

func (a *API) Start() error {
	a.Register(a.r)
	log.Printf("server start on :%s", a.config.TODO_PORT)
	err := http.ListenAndServe(fmt.Sprintf(":%s", a.config.TODO_PORT), a.r)
	if err != nil {
		log.Fatal(fmt.Errorf("error starting server %w", err))
	}

	return nil
}
