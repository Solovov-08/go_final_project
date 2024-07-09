package nextdate

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const FormatDate = "20060102"

// CalculateNextDate вычисляет следующую дату для задачи в соответствии с правилом повторения
func CalculateNextDate(now time.Time, dateStr string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("правило повторения не указано")
	}

	date, err := time.Parse(FormatDate, dateStr)
	if err != nil {
		return "", fmt.Errorf("неверный формат даты: %v", err)
	}

	parts := strings.Fields(repeat)
	rule := parts[0]

	var resultDate time.Time
	switch rule {
	case "":
		if date.Before(now) {
			resultDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		} else {
			resultDate = date
		}
	case "d":
		if len(parts) != 2 {
			return "", errors.New("неверный формат повторения для 'd'")
		}

		daysToInt := make([]int, 0, 7)
		days, err := strconv.Atoi(parts[1])
		if err != nil || days <= 0 || days > 400 {
			return "", errors.New("неверное кол-во дней")
		}
		daysToInt = append(daysToInt, days)

		if daysToInt[0] == 1 {
			resultDate = date.AddDate(0, 0, 1)
		} else {
			resultDate = date.AddDate(0, 0, daysToInt[0])
			for resultDate.Before(now) {
				resultDate = resultDate.AddDate(0, 0, daysToInt[0])
			}
		}
	case "y":
		if len(parts) != 1 {
			return "", errors.New("неверный формат повторения для 'y'")
		}

		resultDate = date.AddDate(1, 0, 0)
		for resultDate.Before(now) {
			resultDate = resultDate.AddDate(1, 0, 0)
		}
	default:
		return "", errors.New("не поддерживаемый формат повторения")
	}

	return resultDate.Format(FormatDate), nil
}
