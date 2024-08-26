package httpsrv

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/EnrikeM/Yandex_final_project_Go/app/internal/calc"
	"github.com/EnrikeM/Yandex_final_project_Go/app/internal/storage"
)

type Response struct {
	MessageKey string
	MessageVal string
	W          http.ResponseWriter
	RespCode   int
}

func WriteResponse(
	messageKey string,
	messageVal string,
	w http.ResponseWriter,
	respCode int) {
	w.WriteHeader(respCode)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(map[string]string{messageKey: messageVal}); err != nil {
		http.Error(w, "error encoding response", respCode)
		return
	}
}

func validate(t *storage.Task) error {
	if t.Date == "" {
		t.Date = time.Now().Format(calc.TimeFormat)
	}

	if _, err := time.Parse(calc.TimeFormat, t.Date); err != nil {
		return fmt.Errorf("field `date` must be in format YYYYMMDD, but provided %w", err)
	}

	if t.Title == "" {
		return fmt.Errorf("field `title` cannot be empty")
	}

	nextDate, err := calc.NextDate(time.Now(), t.Date, t.Repeat)
	if err != nil {
		return fmt.Errorf("couldn't resolve next date: %w", err)
	}

	if t.Date < time.Now().Format(calc.TimeFormat) {
		if t.Repeat == "" {
			now := time.Now().Format(calc.TimeFormat)
			t.Date = now
		} else {
			t.Date = nextDate
		}
	}

	return nil
}
