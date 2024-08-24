package httpsrv

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/EnrikeM/Yandex_final_project_Go/app/internal/apierrors"
)

func (a *API) GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	taskID := r.URL.Query().Get("id")
	if taskID == "" {
		ErrIDNotProvided.Error(w, http.StatusBadRequest)
		return
	}

	task, err := getTask(a.DB, taskID)
	if err != nil {
		rErr := apierrors.New(err.Error())
		rErr.Error(w, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	_ = json.NewEncoder(w).Encode(task)
}

func getTask(db *sql.DB, taskID string) (GetTask, error) { //вынести в utils?
	var task GetTask

	query := "SELECT * FROM scheduler WHERE id = ?"
	row := db.QueryRow(query, taskID)

	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return GetTask{}, err
	}

	return task, nil
}
