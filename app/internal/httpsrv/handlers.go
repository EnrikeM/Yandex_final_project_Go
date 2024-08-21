package httpsrv

import (
	"fmt"
	"net/http"
	"time"

	"github.com/EnrikeM/Yandex_final_project_Go/app/internal/validators"
)

func (a *API) GetHandler(w http.ResponseWriter, r *http.Request) {

	nowStr := r.URL.Query().Get("now")
	now, err := time.Parse("20060102", nowStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(fmt.Sprintf("error parsing time: %s", err)))
		return
	}

	date := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")

	result, err := validators.NextDate(now, date, repeat)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		// _, _ = w.Write([]byte(fmt.Sprintf("error time: %s", err)))
	}

	_, _ = w.Write([]byte(result))
}
