package httpsrv

import (
	"net/http"
	"time"

	"github.com/EnrikeM/Yandex_final_project_Go/app/internal/apierrors"
	"github.com/EnrikeM/Yandex_final_project_Go/app/internal/calc"
)

func (a *API) GetNextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.URL.Query().Get("now")
	now, err := time.Parse(calc.TimeFormat, nowStr)
	if err != nil {
		apierrors.ErrParseTime.Error(w, http.StatusBadRequest)
		return
	}

	date := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")

	result, err := calc.NextDate(now, date, repeat)
	if err != nil {
		rErr := apierrors.New(err.Error())
		rErr.Error(w, http.StatusBadRequest)
		return
	}

	_, _ = w.Write([]byte(result))
}
