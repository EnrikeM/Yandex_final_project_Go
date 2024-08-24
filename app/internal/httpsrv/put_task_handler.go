package httpsrv

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/EnrikeM/Yandex_final_project_Go/app/internal/apierrors"
)

func (a *API) PutTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var task GetTask

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		rErr := apierrors.New(err.Error())
		rErr.Error(w, http.StatusBadRequest)
		return
	}

	if task.ID == "" {
		apierrors.ErrIDNotProvided.Error(w, http.StatusBadRequest)
		return
	}

	var validateTask = Task{
		Date:    task.Date,
		Title:   &task.Title,
		Comment: task.Comment,
		Repeat:  task.Repeat,
	}

	err = validateTask.validate()
	if err != nil {
		rErr := apierrors.New(err.Error())
		rErr.Error(w, http.StatusBadRequest)
		return
	}

	if _, err := getTask(a.DB, task.ID); err != nil {
		apierrors.ErrNoSuchTask.Error(w, http.StatusBadRequest)
		return
	}

	if err = redactTask(a.DB, task); err != nil {
		rErr := apierrors.New(err.Error())
		rErr.Error(w, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	_, _ = w.Write([]byte("{}"))
}

func redactTask(db *sql.DB, task GetTask) error {
	query := `
	UPDATE scheduler 
	SET date = ?, title = ?, comment = ?, repeat = ?, id = ?
	WHERE id = ?;`
	_, err := db.Exec(query, task.Date, &task.Title, task.Comment, task.Repeat, task.ID, task.ID)
	if err != nil {
		return err
	}

	return nil
}
