package httpsrv

import (
	"encoding/json"
	"net/http"

	"github.com/EnrikeM/Yandex_final_project_Go/app/internal/apierrors"
	"github.com/EnrikeM/Yandex_final_project_Go/app/internal/storage"
)

func (a *API) GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	search := r.URL.Query().Get("search")

	tasks, err := storage.GetTasks(a.DB, search)
	if err != nil {
		rErr := apierrors.New(err.Error())
		rErr.Error(w, http.StatusBadRequest) // возможно тут 500 лучше вернуть
		return
	}

	response := map[string][]storage.Task{"tasks": tasks}
	if tasks == nil {
		response["tasks"] = []storage.Task{}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	_ = json.NewEncoder(w).Encode(response)
}
