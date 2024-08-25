package httpsrv

import (
	"net/http"
	"time"

	"github.com/EnrikeM/Yandex_final_project_Go/app/internal/apierrors"
	"github.com/EnrikeM/Yandex_final_project_Go/app/internal/calc"
)

// DeleteTaskHandler godoc
//
//	@Summary		get next date
//	@Description	get next date
//	@Produce		json
//	@Param			now		query		string	true	"current time"
//	@Param			date	query		string	true	"next date"
//	@Param			repeat	query		string	false	"repeat pattern"
//	@Success		200		{object}	string
//	@Failure		400		{object}	nil
//	@Failure		500		{object}	nil
//	@Router			/api/nextdate [get]
//
// .
func (a *API) getNextDateHandler(w http.ResponseWriter, r *http.Request) {
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
