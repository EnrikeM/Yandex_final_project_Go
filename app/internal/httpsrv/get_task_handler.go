package httpsrv

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

type GetTask struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"` //
}

func (a *API) GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	tasks, err := getTasks(a.DB)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(fmt.Sprintf("error getting tasks: %s", err)))
		return
	}

	response := map[string][]GetTask{"tasks": tasks}
	if tasks == nil {
		response["tasks"] = []GetTask{}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err = json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}

}

func getTasks(db *sql.DB) ([]GetTask, error) {
	query := "SELECT * FROM scheduler ORDER BY date DESC LIMIT 10"

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error exectuing query: %v", err)
	}

	defer rows.Close()

	var tasks []GetTask

	for rows.Next() {
		task := GetTask{}
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return tasks, nil
}
