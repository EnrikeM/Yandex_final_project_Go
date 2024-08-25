package httpsrv

import (
	"net/http"

	"github.com/EnrikeM/Yandex_final_project_Go/app/internal/apierrors"
	"github.com/EnrikeM/Yandex_final_project_Go/app/internal/storage"
)

func (a *API) DeleteTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	taskID := r.URL.Query().Get("id")
	if taskID == "" {
		apierrors.ErrIDNotProvided.Error(w, http.StatusBadRequest)
		return
	}

	if _, err := storage.GetTask(a.DB, taskID); err != nil {
		err := apierrors.New(err.Error())
		err.Error(w, http.StatusBadRequest)
		return
	}

	if err := storage.DeleteTask(a.DB, taskID); err != nil {
		err := apierrors.New(err.Error())
		err.Error(w, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	_, _ = w.Write([]byte("{}"))
}
