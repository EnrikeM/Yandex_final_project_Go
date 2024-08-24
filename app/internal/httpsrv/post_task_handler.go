package httpsrv

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/EnrikeM/Yandex_final_project_Go/app/internal/apierrors"
	"github.com/EnrikeM/Yandex_final_project_Go/app/internal/validators"
)

type Task struct {
	Date    string  `json:"date,omitempty"`
	Title   *string `json:"title"`
	Comment string  `json:"comment,omitempty"`
	Repeat  string  `json:"repeat,omitempty"`
}

func (a *API) PostTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var task Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		rErr := apierrors.New(err.Error())
		rErr.Error(w, http.StatusBadRequest)
		return
	}

	err = task.validate()
	if err != nil {
		rErr := apierrors.New(err.Error())
		rErr.Error(w, http.StatusBadRequest)
		return
	}

	if a.DB == nil {
		log.Printf("database connection is nil")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	lastID, err := getLastId(task, a.DB)
	if err != nil {
		log.Printf("error getting last id: %v", err)
		http.Error(w, "error saving task", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	response := map[string]string{"id": lastID}
	if err = json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}
}

func (t *Task) validate() error {
	if t.Date == "" {
		t.Date = time.Now().Format(validators.TimeFormat)
	}

	if _, err := time.Parse(validators.TimeFormat, t.Date); err != nil {
		return fmt.Errorf("field `date` must be in format YYYYMMDD, but provided %w", err)
	}

	if t.Title == nil || *t.Title == "" {
		return fmt.Errorf("field `title` cannot be empty")
	}

	nextDate, err := validators.NextDate(time.Now(), t.Date, t.Repeat)
	if err != nil {
		return fmt.Errorf("couldn't resolve next date: %w", err)
	}

	if t.Date < time.Now().Format(validators.TimeFormat) {
		if t.Repeat == "" {
			now := time.Now().Format(validators.TimeFormat)
			t.Date = now
		} else {
			t.Date = nextDate
		}
	}

	return nil
}

func getLastId(task Task, db *sql.DB) (string, error) {
	query := "INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)"

	result, err := db.Exec(query, task.Date, &task.Title, task.Comment, task.Repeat)
	if err != nil {
		return "", fmt.Errorf("error executing query: %w", err)
	}

	resInt, err := (result.LastInsertId())
	if err != nil {
		return "", fmt.Errorf("error getting last id: %w", err)
	}
	res := strconv.Itoa(int(resInt))

	return res, nil
}
