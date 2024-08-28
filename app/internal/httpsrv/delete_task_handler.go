package httpsrv

import (
	"net/http"

	"github.com/EnrikeM/Yandex_final_project_Go/app/internal/apierrors"
)

// deleteTaskHandler godoc
//
//	@Summary		Delete a task
//	@Description	Delete a task by ID
//	@Produce		json
//	@Param			id	query		int	true	"Task ID"
//	@Success		200	{object}	map[string]string
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/api/task [delete]
//
// .
func (a *API) deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	taskID := r.URL.Query().Get("id")
	if taskID == "" {
		apierrors.ErrIDNotProvided.Error(w, http.StatusBadRequest)
		return
	}

	// if _, err := storage.GetTask(a.DB, taskID); err != nil {
	// 	err := apierrors.New(err.Error())
	// 	err.Error(w, http.StatusBadRequest)
	// 	return
	// }

	if err := a.DB.DeleteTask(taskID); err != nil {
		err := apierrors.New(err.Error())
		err.Error(w, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	_, _ = w.Write([]byte("{}"))
}
