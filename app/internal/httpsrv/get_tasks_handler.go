package httpsrv

import (
	"encoding/json"
	"net/http"

	"github.com/EnrikeM/Yandex_final_project_Go/app/internal/apierrors"
	"github.com/EnrikeM/Yandex_final_project_Go/app/internal/storage"
)

// getTasksHandler godoc
//
//	@Summary		get tasks
//	@Description	get info about all tasks
//	@Produce		json
//	@Param			search	query		string	false	"query"
//	@Success		200		{object}	map[string][]storage.Task
//	@Failure		400		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/api/tasks [get]
//
// .
func (a *API) getTasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	search := r.URL.Query().Get("search")

	tasks, err := a.storage.GetTasks(search)
	if err != nil {
		rErr := apierrors.New(err.Error())
		rErr.Error(w, http.StatusBadRequest)
		return
	}

	response := map[string][]storage.Task{"tasks": tasks}
	if tasks == nil {
		response["tasks"] = []storage.Task{}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	_ = json.NewEncoder(w).Encode(response)
}
