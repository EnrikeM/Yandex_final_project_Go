package httpsrv

import (
	"encoding/json"
	"net/http"

	"github.com/EnrikeM/Yandex_final_project_Go/app/internal/response"
)

// getTaskHandler godoc
//
//	@Summary		get task
//	@Description	get task info by task id
//	@Produce		json
//	@Param			id	query		string	true	"task id"
//	@Success		200	{object}	map[string]string
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/api/task [get]
//
// .
func (a *API) getTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	taskID := r.URL.Query().Get("id")
	if taskID == "" {
		response.ErrIDNotProvided.Error(w, http.StatusBadRequest)
		return
	}

	task, err := a.storage.GetTask(taskID)
	if err != nil {
		rErr := response.New(err.Error())
		rErr.Error(w, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	_ = json.NewEncoder(w).Encode(task)
}
