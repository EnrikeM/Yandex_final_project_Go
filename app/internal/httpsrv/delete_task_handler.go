package httpsrv

import (
	"encoding/json"
	"net/http"
)

func (a *API) DeleteTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	taskID := r.URL.Query().Get("id")
	if taskID == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err := json.NewEncoder(w).Encode(map[string]string{"error": "id not provided"}); err != nil {
			http.Error(w, "error encoding response", http.StatusInternalServerError)
			return
		}
		return
	}

	if _, err := getTask(a.DB, taskID); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err := json.NewEncoder(w).Encode(map[string]string{"error": "нет такого id"}); err != nil {
			http.Error(w, "error encoding response", http.StatusInternalServerError)
			return
		}
		return
	}

	if err := deleteTask(a.DB, taskID); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err := json.NewEncoder(w).Encode(map[string]string{"error": err.Error()}); err != nil {
			http.Error(w, "error encoding response", http.StatusInternalServerError)
			return
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	_, _ = w.Write([]byte("{}"))
}
