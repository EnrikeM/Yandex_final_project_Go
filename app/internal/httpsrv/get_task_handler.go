package httpsrv

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

func (a *API) GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	taskID := r.URL.Query().Get("id")
	if taskID == "" {
		if err := json.NewEncoder(w).Encode(map[string]string{"error": "не указан ID"}); err != nil {
			http.Error(w, "error encoding response", http.StatusInternalServerError)
			return
		}
		return
	}

	task, err := getTask(a.DB, taskID)
	if err != nil {
		if err := json.NewEncoder(w).Encode(map[string]string{"error": "задача не найдена"}); err != nil {
			http.Error(w, "error encoding response", http.StatusInternalServerError)
			return
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err = json.NewEncoder(w).Encode(task); err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}

}

func getTask(db *sql.DB, taskID string) (GetTask, error) {
	var task GetTask

	query := "SELECT * FROM scheduler WHERE id = ?"
	row := db.QueryRow(query, taskID)

	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return GetTask{}, err
	}

	return task, nil
}
