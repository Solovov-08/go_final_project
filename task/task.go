package task

import (
	"go_final_project/nextdate"

	"errors"
	"time"
)

const FormatDate = "20060102"

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// Validate выполняет валидацию задачи на основе определенных правил и возвращает ошибку, если задача не соответствует этим правилам.
func Validate(task *Task) error {
	if task.Title == "" {
		return errors.New("не указан заголовок задачи")
	}

	now := time.Now()
	if task.Date == "" {
		task.Date = now.Format(FormatDate)
	} else {
		date, err := time.Parse(FormatDate, task.Date)
		if err != nil {
			return errors.New("дата представлена в неправильном формате")
		}

		if date.Before(now) {
			if task.Repeat == "" {
				task.Date = now.Format(FormatDate)
			} else {
				nextDate, err := nextdate.CalculateNextDate(now, task.Date, task.Repeat)
				if err != nil {
					return errors.New("ошибка вычисления следующей даты")
				}
				task.Date = nextDate
			}
		}
	}

	return nil
}
