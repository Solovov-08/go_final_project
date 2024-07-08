package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"go_final_project/database"
	"go_final_project/tasks"

	"github.com/jmoiron/sqlx"
)

// MarkTaskDoneHandler обрабатывает запросы на пометку задачи выполненной
func MarkTaskDoneHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		// Получаем идентификатор задачи из параметров запроса
		taskID := r.URL.Query().Get("id")
		if taskID == "" {
			http.Error(w, `{"error": "Не указан идентификатор задачи"}`, http.StatusBadRequest)
			return
		}

		// Получаем задачу из базы данных по идентификатору
		task, err := database.FetchTaskByID(db, taskID)
		if err != nil {
			if err.Error() == "task not found" {
				http.Error(w, `{"error": "Задача не найдена"}`, http.StatusNotFound)
			} else {
				http.Error(w, `{"error": "Ошибка при получении задачи"}`, http.StatusInternalServerError)
			}
			return
		}

		// Если задача не повторяющаяся, удаляем её из базы данных
		if task.Repeat == "" {
			err := database.RemoveTask(db, taskID)
			if err != nil {
				http.Error(w, `{"error": "Ошибка при удалении задачи"}`, http.StatusInternalServerError)
				return
			}
		} else {
			// Если задача повторяющаяся, обновляем её следующей датой выполнения
			now := time.Now()
			nextDate, err := tasks.CalculateNextDate(now, task.Date, task.Repeat)
			if err != nil {
				log.Printf("Ошибка вычисления следующей даты: %v", err)
				http.Error(w, `{"error": "Ошибка вычисления следующей даты"}`, http.StatusInternalServerError)
				return
			}

			// Обновляем задачу в базе данных с новой датой выполнения
			task.Date = nextDate
			err = database.UpdateTask(db, task)
			if err != nil {
				http.Error(w, `{"error": "Ошибка при обновлении задачи"}`, http.StatusInternalServerError)
				return
			}
		}

		// Отправляем пустой JSON в ответ на успешное выполнение
		json.NewEncoder(w).Encode(map[string]interface{}{})
	}
}
