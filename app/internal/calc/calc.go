package calc

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

const TimeFormat = "20060102"

const (
	errInvalidDateFormat  = "param `date` must be in format YYYYMMDD"
	errInvalidNowFormat   = "param `now` must be in format YYYYMMDD"
	errInvalidRepeatValue = "value of timePeriod must be an integer"
	errForbiddenRepeat    = "%s is forbidden, please use `d` or `y`"
	errInvalidDay         = "`m` value must be integer or array with values from -2 to 31 excluding 0"
	errInvalidWeekday     = "`w` value must be integer or array with values from 1 to 7"
	errInvalidQuery       = "params `now` and `date` must be defined in query"
	errForbiddenDUsage    = "param `d` must be followed with a number"
	errForbiddenYUsage    = "`y` must be provided without value"
	errValueDTooBig       = "value of `d` must be less than 400"
	errForbiddenMValue    = "`m` second value must be integer or array with values from 1 to 12"
)

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
	dateFormatted, err := time.Parse(TimeFormat, date)

	if err != nil {
		return "", fmt.Errorf(errInvalidDateFormat)
	}

	if now.IsZero() || date == "" {
		return "", fmt.Errorf(errInvalidQuery)
	}

	if _, err := strconv.Atoi(now.Format(TimeFormat)); err != nil || len(now.Format(TimeFormat)) != 8 {
		return "", fmt.Errorf(errInvalidNowFormat)
	}
	if _, err := strconv.Atoi(date); err != nil || len(date) != 8 {
		return "", fmt.Errorf(errInvalidDateFormat)
	}

	if repeat == "" {
		return "", nil
	}

	repeatSep := strings.Split(repeat, " ")
	repeatMeas := repeatSep[0]

	switch repeatMeas {
	case "d":
		date, err := dayHandler(dateFormatted, now, repeatSep)
		if err != nil {
			return "", err
		}
		return date, nil

	case "y":
		date, err := yearHandler(dateFormatted, now, repeatSep)
		if err != nil {
			return "", err
		}
		return date, nil

	case "w":
		date, err := weekHandler(dateFormatted, now, repeatSep)
		if err != nil {
			return "", err
		}
		return date, nil

	case "m":
		date, err := monthHandler(dateFormatted, now, repeatSep)
		if err != nil {
			return "", err
		}
		return date, nil

	default:
		return "", fmt.Errorf(errForbiddenRepeat, repeat)
	}

}

func dayHandler(dateFormatted, now time.Time, repeatSep []string) (string, error) {
	if len(repeatSep) == 1 {
		return "", fmt.Errorf(errForbiddenDUsage)
	}

	repeatValStr := repeatSep[1]
	repeatVal, err := strconv.Atoi(repeatValStr)
	if err != nil {
		return "", fmt.Errorf(errInvalidRepeatValue)
	}
	if repeatVal > 400 {
		return "", fmt.Errorf(errValueDTooBig)
	}
	if dateFormatted.Sub(now) < 0 {
		dateFormatted = makeDate(dateFormatted, now, 0, 0, repeatVal)
		return dateFormatted.Format(TimeFormat), nil
	}

	return dateFormatted.AddDate(0, 0, repeatVal).Format(TimeFormat), nil
}

func weekHandler(dateFormatted, now time.Time, repeatSep []string) (string, error) {
	if len(repeatSep) == 1 || len(repeatSep) > 2 {
		return "", fmt.Errorf(errInvalidWeekday)
	}
	repeatValStr := repeatSep[1]
	weekdays := strings.Split(repeatValStr, ",")
	for _, weekday := range weekdays {
		weekdayInt, err := strconv.Atoi(weekday)
		if err != nil || weekdayInt < 1 || weekdayInt > 7 {
			return "", fmt.Errorf(errInvalidWeekday)
		}
	}

	dateFormatted = makeDate(dateFormatted, now, 0, 0, 1)

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
}

func monthHandler(dateFormatted, now time.Time, repeatSep []string) (string, error) {
	if len(repeatSep) < 1 {
		return "", fmt.Errorf(errForbiddenMValue)
	}
	repeatVals := repeatSep[1:]
	if len(repeatVals) > 2 {
		return "", fmt.Errorf(errForbiddenMValue)
	}

	days, err := getDays(dateFormatted, repeatVals)
	if err != nil {
		return "", err
	}

	dateFormatted = makeDate(dateFormatted, now, 0, 0, 1)

	if dateFormatted.Day() == days[0] {
		dateFormatted = dateFormatted.AddDate(0, 0, 1)
	}

	val, err := getMonthDate(dateFormatted, days, repeatVals)
	if err != nil {
		return "", err
	}
	return val, nil
}

func yearHandler(dateFormatted, now time.Time, repeatSep []string) (string, error) {
	if len(repeatSep) != 1 {
		return "", fmt.Errorf(errForbiddenYUsage)
	}
	if dateFormatted.Sub(now) < 0 {
		dateFormatted = makeDate(dateFormatted, now, 1, 0, 0)
		return dateFormatted.Format(TimeFormat), nil
	}

	return dateFormatted.AddDate(1, 0, 0).Format(TimeFormat), nil
}

func getDayOfMonth(date time.Time, shift int) int {
	firstOfMonth := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
	LastOfMonth := firstOfMonth.AddDate(0, 1, shift).Day()
	return LastOfMonth
}

func makeDate(dateFormatted, now time.Time, year, month, day int) time.Time {
	for dateFormatted.Sub(now).Hours() < 24 {
		dateFormatted = dateFormatted.AddDate(year, month, day)
	}
	return dateFormatted
}

func getDays(dateFormatted time.Time, repeatVals []string) ([]int, error) {
	daysStr := strings.Split(repeatVals[0], ",")
	var days []int

	for _, val := range daysStr {
		day, err := strconv.Atoi(val)
		if err != nil || day < -2 || day == 0 || day > 31 {
			return nil, fmt.Errorf(errInvalidDay)
		}
		if day == -1 || day == -2 {
			day = getDayOfMonth(dateFormatted, day)
		}
		days = append(days, day)
	}
	sort.Ints(days)
	return days, nil
}

func getMonths(repeatVals []string) ([]int, error) {
	var months []int
	monthsStr := strings.Split(repeatVals[1], ",")

	for _, val := range monthsStr {
		month, err := strconv.Atoi(val)
		if err != nil || month < 1 || month > 12 {
			return nil, fmt.Errorf(errInvalidDay)
		}
		months = append(months, month)

	}
	sort.Ints(months)
	return months, nil
}

func getMonthDate(dateFormatted time.Time, days []int, repeatVals []string) (string, error) {
	if len(repeatVals) == 2 {
		months, err := getMonths(repeatVals)
		if err != nil {
			return "", err
		}

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
