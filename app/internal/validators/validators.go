package validators

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const timeFormat = "20060102"

func NextDate(now time.Time, date string, repeat string) (string, error) {
	var result string
	dateFormatted, err := time.Parse("20060102", date)
	if err != nil {
		return result, fmt.Errorf("param `date` must be in format YYYYMMDD")
	}

	if now.IsZero() || date == "" {
		return result, fmt.Errorf("params `now` and `date` must be defined in query\n")
	}
	if len(now.Format(timeFormat)) != 8 || len(date) != 8 {
		return result, fmt.Errorf("params `now` and `date must be in format YYYYMMDD\n")
	}
	if _, err := strconv.Atoi(now.Format(timeFormat)); err != nil {
		return result, fmt.Errorf("param `now` must be in format YYYYMMDD\n")
	}
	if _, err := strconv.Atoi(date); err != nil {
		return result, fmt.Errorf("param `date` must be in format YYYYMMDD\n")
	}

	if repeat == "" {
		return result, nil //task deleted\n
	} else {
		repeatSep := strings.Split(repeat, " ")
		repeatMeas := repeatSep[0]

		if len(repeatSep) == 1 {
			switch repeatMeas {
			case "d":
				return result, fmt.Errorf("param `d` must be followed with a number \n")
			case "y":
				if dateFormatted.Sub(now) < 0 {
					for dateFormatted.Sub(now) < 0 {
						dateFormatted = dateFormatted.AddDate(1, 0, 0)
					}
					return dateFormatted.Format(timeFormat), nil
				}
				return dateFormatted.AddDate(1, 0, 0).Format(timeFormat), nil
			default:
				return "", fmt.Errorf("%s is forbidden, please use `d`,`w`, `m` or `y", repeat)
			}
		}

		repeatValStr := repeatSep[1]
		repeatVal, err := strconv.Atoi(repeatValStr)
		if err != nil {
			return "", fmt.Errorf("value of timePeriod must be an integer")
		}
		switch repeatMeas {
		case "d":
			if repeatVal > 400 {
				return "", fmt.Errorf("value of `d` must be less than 400")
			}
			if dateFormatted.Sub(now) < 0 {
				for dateFormatted.Sub(now).Hours() < 0 {
					dateFormatted = dateFormatted.AddDate(0, 0, repeatVal)
				}
				return dateFormatted.Format(timeFormat), nil
			}
			return dateFormatted.AddDate(0, 0, repeatVal).Format(timeFormat), nil

		case "y":
			return "", fmt.Errorf("`y` must be provided without value")
		}

	}

	return "", nil
}
