package httpsrv

import (
	"encoding/json"
	"net/http"

	"github.com/EnrikeM/Yandex_final_project_Go/app/internal/apierrors"
	"github.com/EnrikeM/Yandex_final_project_Go/app/internal/storage"
)

func (a *API) GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	taskID := r.URL.Query().Get("id")
	if taskID == "" {
		apierrors.ErrIDNotProvided.Error(w, http.StatusBadRequest)
		return
	}

	task, err := storage.GetTask(a.DB, taskID)
	if err != nil {
		rErr := apierrors.New(err.Error())
		rErr.Error(w, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	_ = json.NewEncoder(w).Encode(task)
}

// func getTask(db *sql.DB, taskID string) (storage.GetTask, error) { //вынести в utils?
// 	var task storage.GetTask

// 	query := "SELECT * FROM scheduler WHERE id = ?"
// 	row := db.QueryRow(query, taskID)

// 	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
// 	if err != nil {
// 		return storage.GetTask{}, err
// 	}

// 	return task, nil
// }
