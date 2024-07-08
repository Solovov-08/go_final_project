package database

import (
	"database/sql"
	"errors"
	"go_final_project/model"
	"go_final_project/tasks"
	"log"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

// FetchTaskByID берёт задачу из базы данных по идентификатору
func FetchTaskByID(db *sqlx.DB, id string) (model.Task, error) {
	var task model.Task
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
	err := db.QueryRow(query, id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return task, errors.New("задача не найдена")
		}
		log.Printf("Ошибка выполнения запроса: %v", err)
		return task, errors.New("ошибка выполнения запроса")
	}
	return task, nil
}

// FetchTasks выполняет поиск задач в базе данных с учетом поискового запроса и ограничения на количество возвращаемых задач.
func FetchTasks(db *sqlx.DB, search string, limit int) ([]model.Task, error) {
	var rows *sqlx.Rows
	var err error

	// Убираем пробелы в строке
	search = strings.TrimSpace(search)
	if search != "" {
		// Проверяем, является ли строка датой в формате "02.01.2006"
		if searchDate, err := time.Parse("02.01.2006", search); err == nil {
			searchDateStr := searchDate.Format(tasks.FormatDate)
			rows, err = db.NamedQuery(`SELECT id, date, title, comment, repeat FROM scheduler WHERE date = :searchDate ORDER BY date LIMIT :limit`, map[string]interface{}{
				"searchDate": searchDateStr,
				"limit":      limit,
			})
			if err != nil {
				log.Printf("Ошибка выполнения запроса с датой: %v", err)
				return nil, err
			}
		} else {
			rows, err = db.NamedQuery(`SELECT id, date, title, comment, repeat FROM scheduler WHERE LOWER(title) LIKE :search OR LOWER(comment) LIKE :search ORDER BY date LIMIT :limit`, map[string]interface{}{
				"search": "%" + search + "%",
				"limit":  limit,
			})
			if err != nil {
				log.Printf("Ошибка выполнения запроса с поиском: %v", err)
				return nil, err
			}
		}
	} else {
		// Если строка пустая, просто выбираем все задачи
		rows, err = db.NamedQuery(`SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT :limit`, map[string]interface{}{
			"limit": limit,
		})
		if err != nil {
			log.Printf("Ошибка выполнения запроса без поиска: %v", err)
			return nil, err
		}
	}

	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("Ошибка закрытия rows: %v", err)
		}
	}()

	var tasks []model.Task
	for rows.Next() {
		var task model.Task
		err := rows.StructScan(&task)
		if err != nil {
			log.Printf("Ошибка сканирования строки: %v", err)
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Ошибка чтения строк: %v", err)
		return nil, err
	}

	// Возвращаем пустой массив задач, если ни одной задачи не найдено
	if tasks == nil {
		tasks = []model.Task{}
	}

	return tasks, nil
}
