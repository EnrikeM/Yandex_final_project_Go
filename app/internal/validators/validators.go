package validators

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"
)

// TODO: Refactor
const TimeFormat = "20060102"

var weekDay = map[string]string{
	"1": "Monday",
	"2": "Tuesday",
	"3": "Wednesday",
	"4": "Thursday",
	"5": "Friday",
	"6": "Saturday",
	"7": "Sunday",
}

var monthsMap = map[int]string{
	1:  "January",
	2:  "February",
	3:  "Marth",
	4:  "April",
	5:  "May",
	6:  "June",
	7:  "July",
	8:  "August",
	9:  "September",
	10: "October",
	11: "November",
	12: "December",
}

func NextDate(now time.Time, date string, repeat string) (string, error) {
	log.Println("NEXT DATE")
	dateFormatted, err := time.Parse(TimeFormat, date)

	log.Printf("Now: %v, Date: %s", now, date)
	if err != nil {
		return "", fmt.Errorf("param `date` must be in format YYYYMMDD")
	}
	log.Println("After dateFormatted")

	if now.IsZero() || date == "" {
		return "", fmt.Errorf("params `now` and `date` must be defined in query\n")
	}

	log.Println("After first ")

	if _, err := strconv.Atoi(now.Format(TimeFormat)); err != nil || len(now.Format(TimeFormat)) != 8 {
		return "", fmt.Errorf("param `now` must be in format YYYYMMDD\n")
	}
	if _, err := strconv.Atoi(date); err != nil || len(date) != 8 {
		return "", fmt.Errorf("param `date` must be in format YYYYMMDD\n")
	}

	log.Println("After all init things")

	if repeat == "" {
		return "", nil //task deleted\n
	}
	repeatSep := strings.Split(repeat, " ")
	repeatMeas := repeatSep[0]

	if len(repeatSep) == 1 {
		log.Println("solo")
		switch repeatMeas {
		case "d":
			return "", fmt.Errorf("param `d` must be followed with a number \n")
		case "y":
			if dateFormatted.Sub(now) < 0 {
				for dateFormatted.Sub(now) < 0 {
					dateFormatted = dateFormatted.AddDate(1, 0, 0)
				}
				return dateFormatted.Format(TimeFormat), nil
			}
			return dateFormatted.AddDate(1, 0, 0).Format(TimeFormat), nil
		default:
			return "", fmt.Errorf("%s is forbidden, please use `d` or `y", repeat)
		}
	}

	repeatValStr := repeatSep[1]

	switch repeatMeas {
	case "d":
		repeatVal, err := strconv.Atoi(repeatValStr)
		if err != nil {
			return "", fmt.Errorf("value of timePeriod must be an integer")
		}

		if repeatVal > 400 {
			return "", fmt.Errorf("value of `d` must be less than 400")
		}
		if dateFormatted.Sub(now) < 0 {
			for dateFormatted.Sub(now) < 0 {
				dateFormatted = dateFormatted.AddDate(0, 0, repeatVal)
			} // Возможно, стоит убрать if
			return dateFormatted.Format(TimeFormat), nil
		}
		return dateFormatted.AddDate(0, 0, repeatVal).Format(TimeFormat), nil

	case "y":
		return "", fmt.Errorf("`y` must be provided without value")

	case "w":
		if len(repeatSep) > 2 {
			return "", fmt.Errorf("`w` value must be integer or array with values from 1 to 7")
		}
		weekdays := strings.Split(repeatValStr, ",") //[4,5]
		for _, weekday := range weekdays {
			weekdayInt, err := strconv.Atoi(weekday)
			if err != nil || weekdayInt < 1 || weekdayInt > 7 {
				return "", fmt.Errorf("`w` value must be integer or array with values from 1 to 7")
			}
		}

		if dateFormatted.Sub(now).Hours() < 24 {
			for dateFormatted.Sub(now).Hours() < 24 {
				dateFormatted = dateFormatted.AddDate(0, 0, 1)
			}
		} // убрать иф?

		if dateFormatted.String() == weekDay[weekdays[0]] {
			dateFormatted = dateFormatted.AddDate(0, 0, 1)
		}

		for {
			for _, weekday := range weekdays {
				if dateFormatted.Weekday().String() == weekDay[weekday] {
					return dateFormatted.Format(TimeFormat), nil
				}
			}
			dateFormatted = dateFormatted.AddDate(0, 0, 1)
		}

	case "m":
		repeatVals := repeatSep[1:] //  m ["4,5", "10,11"]
		if len(repeatVals) > 2 {
			log.Println("double")
			return "", fmt.Errorf("`m` second value must be integer or array with values from 1 to 12")
		}

		daysStr := strings.Split(repeatVals[0], ",")
		var days []int

		for _, val := range daysStr {
			day, err := strconv.Atoi(val)
			if err != nil || day < -2 || day == 0 || day > 31 {
				return "", fmt.Errorf("`m` value must be integer or array with values from -2 to 31")
			}
			if day == -1 || day == -2 {
				day = getDayOfMonth(dateFormatted, day)
			}
			days = append(days, day)
		}
		sort.Ints(days)

		if dateFormatted.Sub(now).Hours() < 24 {
			for dateFormatted.Sub(now).Hours() < 24 {
				dateFormatted = dateFormatted.AddDate(0, 0, 1)
			}
		}

		if dateFormatted.Day() == days[0] {
			dateFormatted = dateFormatted.AddDate(0, 0, 1)
		}

		if len(repeatVals) == 2 {
			var months []int
			monthsStr := strings.Split(repeatVals[1], ",")

			for _, val := range monthsStr {
				month, err := strconv.Atoi(val)
				if err != nil || month < 1 || month > 12 {
					return "", fmt.Errorf("`m` value must be integer or array with values from -2 to 31")
				}
				months = append(months, month)

			}
			sort.Ints(months)

		dayStart:
			for {
				for _, day := range days {
					if dateFormatted.Day() == day {
						break dayStart
					}
				}
				dateFormatted = dateFormatted.AddDate(0, 0, 1)
			}

			for {
				for _, month := range months {
					if dateFormatted.Month().String() == monthsMap[month] {
						return dateFormatted.Format(TimeFormat), nil
					}
				}
				dateFormatted = dateFormatted.AddDate(0, 1, 0)
			}

		}

		for {
			for _, day := range days {
				if dateFormatted.Day() == day {
					return dateFormatted.Format(TimeFormat), nil
				}
			}
			dateFormatted = dateFormatted.AddDate(0, 0, 1)
		}

	}

	log.Println("NEXT DATE END")
	return "", nil
}

func getDayOfMonth(date time.Time, shift int) int {
	firstOfMonth := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
	LastOfMonth := firstOfMonth.AddDate(0, 1, shift).Day()
	return LastOfMonth
}
