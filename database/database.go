package database

import (
	"database/sql"
	"errors"
	"go_final_project/nextdate"
	"go_final_project/task"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

const dbFileName = "scheduler.db"

type Storage struct {
	db *sqlx.DB
}

func NewStorage(db *sqlx.DB) *Storage {
	return &Storage{db: db}
}

// InitializeDatabase проверяет существование файла базы данных и создает таблицу, если необходимо
func InitializeDatabase() (*Storage, error) {
	if _, err := os.Stat(dbFileName); os.IsNotExist(err) {
		file, err := os.Create(dbFileName)
		if err != nil {
			return nil, err
		}
		file.Close()
	}

	db, err := sqlx.Connect("sqlite3", dbFileName)
	if err != nil {
		return nil, err
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS scheduler (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
		date CHAR(8) NOT NULL DEFAULT "",
		title VARCHAR(128) NOT NULL DEFAULT "",
		comment TEXT NOT NULL DEFAULT "",
		repeat VARCHAR(128) NOT NULL DEFAULT ""
    );`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return nil, err
	}

	// Создаем индекс по полю date для сортировки задач по дате
	createIndexSQL := `CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler(date);`
	_, err = db.Exec(createIndexSQL)
	if err != nil {
		return nil, err
	}

	storage := NewStorage(db)
	return storage, nil
}

// Close закрывает соединение с базой данных
func (s *Storage) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

// InsertTask добавляет новую задачу в базу данных.
func (s *Storage) InsertTask(task task.Task) (int64, error) {
	if task.Date == "" {
		task.Date = time.Now().Format(nextdate.FormatDate)
	} else {
		_, err := time.Parse(nextdate.FormatDate, task.Date)
		if err != nil {
			return 0, errors.New("дата представлена в неправильном формате")
		}
	}

	now := time.Now()
	if task.Date < now.Format(nextdate.FormatDate) {
		if task.Repeat == "" {
			task.Date = now.Format(nextdate.FormatDate)
		} else {
			nextDate, err := nextdate.CalculateNextDate(now, task.Date, task.Repeat)
			if err != nil {
				return 0, err
			}
			task.Date = nextDate
		}
	}

	result, err := s.db.Exec(
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
func (s *Storage) RemoveTask(id string) error {
	deleteQuery := `DELETE FROM scheduler WHERE id = ?`
	res, err := s.db.Exec(deleteQuery, id)
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

// UpdateTask обновляет задачу в базе данных
func (s *Storage) UpdateTask(task task.Task) error {
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	res, err := s.db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
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

// GetTask берёт задачу из базы данных по идентификатору
func (s *Storage) GetTask(id string) (task.Task, error) {
	var task task.Task
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
	err := s.db.QueryRow(query, id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return task, errors.New("задача не найдена")
		}
		log.Printf("Ошибка выполнения запроса: %v", err)
		return task, errors.New("ошибка выполнения запроса")
	}
	return task, nil
}

// GetTasks выполняет поиск задач в базе данных с учетом поискового запроса и ограничения на количество возвращаемых задач.
func (s *Storage) GetTasks(search string, limit int) ([]task.Task, error) {
	var rows *sqlx.Rows
	var err error

	// Убираем пробелы в строке
	search = strings.TrimSpace(search)
	if search != "" {
		// Проверяем, является ли строка датой в формате "02.01.2006"
		if searchDate, err := time.Parse("02.01.2006", search); err == nil {
			searchDateStr := searchDate.Format(nextdate.FormatDate)
			rows, err = s.db.NamedQuery(`SELECT id, date, title, comment, repeat FROM scheduler WHERE date = :searchDate ORDER BY date LIMIT :limit`, map[string]interface{}{
				"searchDate": searchDateStr,
				"limit":      limit,
			})
			if err != nil {
				log.Printf("Ошибка выполнения запроса с датой: %v", err)
				return nil, err
			}
		} else {
			rows, err = s.db.NamedQuery(`SELECT id, date, title, comment, repeat FROM scheduler WHERE LOWER(title) LIKE :search OR LOWER(comment) LIKE :search ORDER BY date LIMIT :limit`, map[string]interface{}{
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
		rows, err = s.db.NamedQuery(`SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT :limit`, map[string]interface{}{
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

	var tasks []task.Task
	for rows.Next() {
		var task task.Task
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
		tasks = []task.Task{}
	}

	return tasks, nil
}
