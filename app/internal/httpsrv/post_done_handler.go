package httpsrv

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/EnrikeM/Yandex_final_project_Go/app/internal/validators"
)

func (a *API) PostDoneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	taskID := r.URL.Query().Get("id")
	if taskID == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err := json.NewEncoder(w).Encode(map[string]string{"error": "id cannot be null"}); err != nil {
			http.Error(w, "error encoding response", http.StatusInternalServerError)
			return
		}
		return
	}
	//Вынести errors.go, куда добавить конструктор ошибок и метод String, куда будет писаться message

	task, err := getTask(a.DB, taskID)
	log.Println(task)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err = json.NewEncoder(w).Encode(map[string]string{"error": "error searching for task"}); err != nil {
			http.Error(w, "error encoding response", http.StatusInternalServerError)
			return
		}
		return
	}

	if task.Repeat == "" {
		if err := deleteTask(a.DB, taskID); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			if err = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()}); err != nil {
				http.Error(w, "error encoding response", http.StatusInternalServerError)
				return
			}
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		_, _ = w.Write([]byte("{}"))
		return
	}

	task.Date, err = validators.NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err = json.NewEncoder(w).Encode(map[string]string{"error": "error getting next date"}); err != nil {
			http.Error(w, "error encoding response", http.StatusInternalServerError)
			return
		}
		return
	}

	if err = redactTask(a.DB, task); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err = json.NewEncoder(w).Encode(map[string]string{"error": "error redacting task"}); err != nil {
			http.Error(w, "error encoding response", http.StatusInternalServerError)
			return
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	_, _ = w.Write([]byte("{}"))
}

func deleteTask(db *sql.DB, taskID string) error {
	query := "DELETE FROM scheduler WHERE id = ?"
	_, err := db.Exec(query, taskID)
	if err != nil {
		return err
	}
	return nil
}
