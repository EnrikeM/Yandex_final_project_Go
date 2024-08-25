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

// func getTasks(db *sql.DB, search string) ([]storage.GetTask, error) { //вынести в utils?
// 	var query string

// 	if search == "" {
// 		query = "SELECT * FROM scheduler ORDER BY date DESC LIMIT 10"
// 	}
// 	if search != "" {
// 		searchDate, err := time.Parse("02.01.2006", search)
// 		if err == nil {
// 			search = searchDate.Format("20060102")
// 			query = "SELECT * FROM scheduler WHERE date LIKE ? ORDER BY date DESC LIMIT 10"
// 		} else {
// 			query = "SELECT * FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT 10"
// 		}
// 	}

// 	rows, err := db.Query(query, fmt.Sprintf("%%%s%%", search), fmt.Sprintf("%%%s%%", search))
// 	if err != nil {
// 		return nil, fmt.Errorf("error exectuing query: %v", err)
// 	}

// 	defer rows.Close()

// 	var tasks []storage.GetTask

// 	for rows.Next() {
// 		task := storage.GetTask{}
// 		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
// 		if err != nil {
// 			return nil, err
// 		}
// 		tasks = append(tasks, task)
// 	}

// 	if rows.Err() != nil {
// 		return nil, err
// 	}

// 	return tasks, nil
// }
