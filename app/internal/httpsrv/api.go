package httpsrv

import (
	"fmt"
	"log"
	"net/http"

	"github.com/EnrikeM/Yandex_final_project_Go/app/internal/config"
	"github.com/EnrikeM/Yandex_final_project_Go/app/internal/storage"
	"github.com/go-chi/chi"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/EnrikeM/Yandex_final_project_Go/docs"
)

type API struct {
	storage storage.Scheduler
	config  config.Config
	router  chi.Router
}

func NewAPI(config config.Config, storage storage.Scheduler) *API {
	api := &API{
		router:  chi.NewRouter(),
		config:  config,
		storage: storage,
	}

	api.register(api.router)
	return api
}

func (a *API) register(r chi.Router) {

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
		r.Use(a.auth)
		r.Get("/", a.getTasksHandler)
	})

	r.Get("/api/nextdate", a.getNextDateHandler)

	r.Get("/swagger/*", httpSwagger.WrapHandler)
	r.Handle("/*", http.FileServer(http.Dir(a.config.WEB_DIR)))
}

func (a *API) Start() error {
	log.Printf("server start on :%s", a.config.TODO_PORT)
	err := http.ListenAndServe(fmt.Sprintf(":%s", a.config.TODO_PORT), a.router)
	if err != nil {
		return fmt.Errorf("error starting server %w", err)
	}

	return nil
}
