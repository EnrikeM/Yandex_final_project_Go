package httpsrv

import (
	"fmt"
	"time"

	"github.com/EnrikeM/Yandex_final_project_Go/app/internal/calc"
	"github.com/EnrikeM/Yandex_final_project_Go/app/internal/storage"
)

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
