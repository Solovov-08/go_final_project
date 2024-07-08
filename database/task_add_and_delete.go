package database

import (
	"errors"
	"go_final_project/model"
	"go_final_project/tasks"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

// InsertTask добавляет новую задачу в базу данных.
func InsertTask(db *sqlx.DB, task model.Task) (int64, error) {
	if task.Title == "" {
		return 0, errors.New("не указан заголовок задачи")
	}

	if task.Date == "" {
		task.Date = time.Now().Format(tasks.FormatDate)
	} else {
		_, err := time.Parse(tasks.FormatDate, task.Date)
		if err != nil {
			return 0, errors.New("дата представлена в неправильном формате")
		}
	}

	now := time.Now()
	if task.Date < now.Format(tasks.FormatDate) {
		if task.Repeat == "" {
			task.Date = now.Format(tasks.FormatDate)
		} else {
			nextDate, err := tasks.CalculateNextDate(now, task.Date, task.Repeat)
			if err != nil {
				return 0, err
			}
			task.Date = nextDate
		}
	}

	result, err := db.Exec(
		`INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`,
		task.Date, task.Title, task.Comment, task.Repeat,
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// RemoveTask удаляет задачу из базы данных по идентификатору
func RemoveTask(db *sqlx.DB, id string) error {
	deleteQuery := `DELETE FROM scheduler WHERE id = ?`
	res, err := db.Exec(deleteQuery, id)
	if err != nil {
		log.Printf("Ошибка выполнения запроса: %v", err)
		return errors.New("ошибка выполнения запроса")
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Printf("Ошибка получения результата запроса: %v", err)
		return errors.New("ошибка получения результата запроса")
	}

	if rowsAffected == 0 {
		return errors.New("задача не найдена")
	}

	return nil
}
