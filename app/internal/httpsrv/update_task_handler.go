package httpsrv

import (
	"encoding/json"
	"net/http"

	"github.com/EnrikeM/Yandex_final_project_Go/app/internal/response"
	"github.com/EnrikeM/Yandex_final_project_Go/app/internal/storage"
)

// updateTaskHandler godoc
//
//	@Summary		update task
//	@Description	update task with task attributes
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	map[string]string
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/api/task [put]
//
// .
func (a *API) updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var task storage.Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		rErr := response.New(err.Error())
		rErr.Error(w, http.StatusBadRequest)
		return
	}

	if task.ID == "" {
		response.ErrIDNotProvided.Error(w, http.StatusBadRequest)
		return
	}

	err = task.Validate()
	if err != nil {
		rErr := response.New(err.Error())
		rErr.Error(w, http.StatusBadRequest)
		return
	}

	if err = a.storage.Update(task); err != nil {
		rErr := response.New(err.Error())
		rErr.Error(w, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	_, _ = w.Write([]byte("{}"))
}
