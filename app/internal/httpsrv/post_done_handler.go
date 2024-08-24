package httpsrv

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/EnrikeM/Yandex_final_project_Go/app/internal/apierrors"
	"github.com/EnrikeM/Yandex_final_project_Go/app/internal/validators"
)

func (a *API) PostDoneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	taskID := r.URL.Query().Get("id")
	if taskID == "" {
		err := apierrors.ErrIDNotProvided
		err.Error(w, http.StatusBadRequest)
		return
	}

	task, err := getTask(a.DB, taskID)
	if err != nil {
		rErr := apierrors.New(err.Error())
		rErr.Error(w, http.StatusBadRequest)
		return
	}

	if task.Repeat == "" {
		if err := deleteTask(a.DB, taskID); err != nil {
			rErr := apierrors.New(err.Error())
			rErr.Error(w, http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		_, _ = w.Write([]byte("{}"))
		return
	}

	task.Date, err = validators.NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		rErr := apierrors.New(err.Error())
		rErr.Error(w, http.StatusBadRequest) // возможно тут 500 лучше вернуть
		return
	}

	if err = redactTask(a.DB, task); err != nil {
		rErr := apierrors.New(err.Error())
		rErr.Error(w, http.StatusBadRequest) // возможно тут тоже 500 лучше вернуть
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
