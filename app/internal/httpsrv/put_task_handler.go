package httpsrv

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

func (a *API) PutTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var task GetTask

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		if err := json.NewEncoder(w).Encode(map[string]string{"error": err.Error()}); err != nil {
			http.Error(w, "error encoding response", http.StatusInternalServerError)
			return
		}
		return
	}

	if task.ID == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err = json.NewEncoder(w).Encode(map[string]string{"error": "не указан ID"}); err != nil {
			http.Error(w, "error encoding response", http.StatusInternalServerError)
			return
		}
		return
	}

	var validateTask = Task{
		Date:    task.Date,
		Title:   &task.Title,
		Comment: task.Comment,
		Repeat:  task.Repeat,
	}

	err = validateTask.validate()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err := json.NewEncoder(w).Encode(map[string]string{"error": "не валидная задача"}); err != nil {
			http.Error(w, "error encoding response", http.StatusInternalServerError)
			return
		}
		return
	}

	if _, err := getTask(a.DB, task.ID); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err := json.NewEncoder(w).Encode(map[string]string{"error": "нет такого id"}); err != nil {
			http.Error(w, "error encoding response", http.StatusInternalServerError)
			return
		}
		return
	}

	if err = redactTask(a.DB, task); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err := json.NewEncoder(w).Encode(map[string]string{"error": err.Error()}); err != nil {
			return // Вынести
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	_, _ = w.Write([]byte("{}"))

}

func redactTask(db *sql.DB, task GetTask) error {

	query := `
	UPDATE scheduler 
	SET date = ?, title = ?, comment = ?, repeat = ?, id = ?
	WHERE id = ?;`
	_, err := db.Exec(query, task.Date, &task.Title, task.Comment, task.Repeat, task.ID, task.ID)
	if err != nil {
		return err
	}

	return nil
}
