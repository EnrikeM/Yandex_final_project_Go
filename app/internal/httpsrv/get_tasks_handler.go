package httpsrv

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type GetTask struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func (a *API) GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	search := r.URL.Query().Get("search")

	tasks, err := getTasks(a.DB, search)
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

func getTasks(db *sql.DB, search string) ([]GetTask, error) { //вынести в utils?
	var query string

	if search == "" {
		query = "SELECT * FROM scheduler ORDER BY date DESC LIMIT 10"
	}
	if search != "" {
		searchDate, err := time.Parse("02.01.2006", search)
		if err == nil {
			search = searchDate.Format("20060102")
			query = "SELECT * FROM scheduler WHERE date LIKE ? ORDER BY date DESC LIMIT 10"
		} else {

			query = "SELECT * FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT 10"
		}
	}

	rows, err := db.Query(query, fmt.Sprintf("%%%s%%", search), fmt.Sprintf("%%%s%%", search))
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
