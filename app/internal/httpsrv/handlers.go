package httpsrv

import (
	"net/http"
)

func (a *API) GetHandler(w http.ResponseWriter, r *http.Request) {
	now := r.URL.Query().Get("now")
	date := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")

	if repeat == "" || now == "" || date == "" {
		return
	}

	// // dateSep := strings.Split(date, " ")
	// dateSep := []byte(date)

	// dateMeas := dateSep[0]
	// dateVal := dateSep[1]
	// repeat := strings.Split(repeat, " ")

}

// "api/nextdate?now=20240126&date=20240126&repeat=y"
