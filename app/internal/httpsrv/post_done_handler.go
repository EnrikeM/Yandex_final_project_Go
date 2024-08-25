package httpsrv

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/EnrikeM/Yandex_final_project_Go/app/internal/apierrors"
	"github.com/EnrikeM/Yandex_final_project_Go/app/internal/calc"
	"github.com/EnrikeM/Yandex_final_project_Go/app/internal/storage"
)

// postDoneHandler godoc
//
//	@Summary		mark  task as done
//	@Description	mark  task as done by task id
//	@Produce		json
//	@Param			id	query		string	true	"id"
//	@Success		200	{object}	nil
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/api/task/done [post]
//
// .
func (a *API) postDoneHandler(w http.ResponseWriter, r *http.Request) {
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

	task, err := storage.GetTask(a.DB, taskID)
	if err != nil {
		rErr := apierrors.New(err.Error())
		rErr.Error(w, http.StatusBadRequest)
		return
	}

	if task.Repeat == "" {
		if err := storage.DeleteTask(a.DB, taskID); err != nil {
			rErr := apierrors.New(err.Error())
			rErr.Error(w, http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		_, _ = w.Write([]byte("{}"))
		return
	}

	task.Date, err = calc.NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		rErr := apierrors.New(err.Error())
		rErr.Error(w, http.StatusBadRequest) // возможно тут 500 лучше вернуть
		return
	}

	if err = storage.RedactTask(a.DB, task); err != nil {
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
