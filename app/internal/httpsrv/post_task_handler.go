package httpsrv

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/EnrikeM/Yandex_final_project_Go/app/internal/response"
	"github.com/EnrikeM/Yandex_final_project_Go/app/internal/storage"
)

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
		rErr := response.New(err.Error())
		rErr.Error(w, http.StatusBadRequest)
		return
	}

	err = task.Validate()
	if err != nil {
		rErr := response.New(err.Error())
		rErr.Error(w, http.StatusBadRequest)
		return
	}

	lastID, err := a.storage.Add(task)
	if err != nil {
		log.Printf("error getting last id: %v", err)
		http.Error(w, "error saving task", http.StatusInternalServerError)
		return
	}

	response.WriteResponse("id", lastID, w, http.StatusOK)

}
