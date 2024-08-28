package httpsrv

import (
	"database/sql"
	"encoding/json"
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
	DB     storage.Scheduler
	config config.Config
	router chi.Router
}

func NewAPI(db *sql.DB, config config.Config) *API {
	api := &API{
		router: chi.NewRouter(),
		config: config,
		DB:     storage.New(db, config),
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

type Response struct {
	MessageKey string
	MessageVal string
	W          http.ResponseWriter
	RespCode   int
}

func WriteResponse(
	messageKey string,
	messageVal string,
	w http.ResponseWriter,
	respCode int) {
	w.WriteHeader(respCode)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(map[string]string{messageKey: messageVal}); err != nil {
		http.Error(w, "error encoding response", respCode)
		return
	}
}
