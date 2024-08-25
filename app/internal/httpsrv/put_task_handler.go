package httpsrv

import (
	"encoding/json"
	"net/http"

	"github.com/EnrikeM/Yandex_final_project_Go/app/internal/apierrors"
	"github.com/EnrikeM/Yandex_final_project_Go/app/internal/storage"
)

func (a *API) PutTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var task storage.Task

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

	if _, err := storage.GetTask(a.DB, task.ID); err != nil {
		apierrors.ErrNoSuchTask.Error(w, http.StatusBadRequest)
		return
	}

	if err = storage.RedactTask(a.DB, task); err != nil {
		rErr := apierrors.New(err.Error())
		rErr.Error(w, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	_, _ = w.Write([]byte("{}"))
}
