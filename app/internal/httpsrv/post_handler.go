package httpsrv

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/EnrikeM/Yandex_final_project_Go/app/internal/validators"
)

type Task struct {
	Date    *string `json:"date"`
	Title   *string `json:"title"`
	Comment string  `json:"comment,omitempty"`
	Repeat  *string `json:"repeat,omitempty"`
}

func (t *Task) validate() error {
	nextDate, err := validators.NextDate(time.Now(), *t.Date, *t.Repeat) // тут может быть проблема
	if err != nil {
		return (fmt.Errorf("error resolving next date for time %s", nextDate))
	}

	if t.Title == nil {
		return (fmt.Errorf("field `title` cannot be empty"))
	}
	if t.Date == nil || *t.Date == "" {
		*t.Date = (time.Now().Format(validators.TimeFormat))
	}
	if *t.Date < (time.Now().Format(validators.TimeFormat)) {
		if t.Repeat == nil || *t.Repeat == "" {
			*t.Date = (time.Now().Format(validators.TimeFormat))
		}
		// nextDate, err := validators.NextDate(time.Now(), *t.Date, *t.Repeat) // тут может быть проблема
		// if err != nil {
		// 	return (fmt.Errorf("error resolving next date for time %s", nextDate))
		// }
	}

	if _, err := time.Parse("20060102", *t.Date); err != nil {
		return (fmt.Errorf("field `date` must be in format YYYYMMDD"))
	}
	return nil
}

func (a *API) PostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var task Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		log.Fatal(fmt.Errorf("error decoding request json %w", err))
	}

	err = task.validate()
	if err != nil {
		log.Fatal(fmt.Errorf("error validating json %w", err))
	}

}
