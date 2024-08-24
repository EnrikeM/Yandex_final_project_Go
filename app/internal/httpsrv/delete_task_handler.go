package httpsrv

import (
	"net/http"

	"github.com/EnrikeM/Yandex_final_project_Go/app/internal/apierrors"
)

var (
	ErrIDNotProvided = apierrors.New("id not provided")
)

func (a *API) DeleteTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	taskID := r.URL.Query().Get("id")
	if taskID == "" {
		ErrIDNotProvided.Error(w, http.StatusBadRequest)
		return
	}

	if _, err := getTask(a.DB, taskID); err != nil {
		err := apierrors.New(err.Error())
		err.Error(w, http.StatusBadRequest)
		return
	}

	if err := deleteTask(a.DB, taskID); err != nil {
		err := apierrors.New(err.Error())
		err.Error(w, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	_, _ = w.Write([]byte("{}"))
}
