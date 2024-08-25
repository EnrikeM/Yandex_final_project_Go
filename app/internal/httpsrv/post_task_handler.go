package httpsrv

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/EnrikeM/Yandex_final_project_Go/app/internal/apierrors"
	"github.com/EnrikeM/Yandex_final_project_Go/app/internal/storage"
)

type Task struct {
	Date    string  `json:"date,omitempty"`
	Title   *string `json:"title"`
	Comment string  `json:"comment,omitempty"`
	Repeat  string  `json:"repeat,omitempty"`
	// id      string  `json:"-"`
}

// postTaskHandler godoc
//
//	@Summary		post task
//	@Description	post task with task attributes
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	map[string]string
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/api/task [post]
//
// .
func (a *API) postTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var task storage.Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		rErr := apierrors.New(err.Error())
		rErr.Error(w, http.StatusBadRequest)
		return
	}

	err = validate(&task)
	if err != nil {
		rErr := apierrors.New(err.Error())
		rErr.Error(w, http.StatusBadRequest)
		return
	}

	if a.DB == nil {
		log.Printf("database connection is nil")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	lastID, err := storage.GetLastId(task, a.DB)
	if err != nil {
		log.Printf("error getting last id: %v", err)
		http.Error(w, "error saving task", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	response := map[string]string{"id": lastID}
	if err = json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}
}
