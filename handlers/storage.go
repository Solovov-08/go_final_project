package handlers

import (
	"encoding/json"
	"go_final_project/database"
	"go_final_project/nextdate"
	"go_final_project/task"
	"log"
	"net/http"
	"time"
)

const TaskLimit = 50

// MarkTaskDoneHandler обрабатывает запросы на пометку задачи выполненной
func MarkTaskDoneHandler(storage *database.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		// Получаем идентификатор задачи из параметров запроса
		taskID := r.URL.Query().Get("id")
		if taskID == "" {
			http.Error(w, `{"error": "Не указан идентификатор задачи"}`, http.StatusBadRequest)
			return
		}

		// Получаем задачу из базы данных по идентификатору
		task, err := storage.GetTask(taskID)
		if err != nil {
			if err.Error() == "задача не найдена" {
				http.Error(w, `{"error": "Задача не найдена"}`, http.StatusNotFound)
			} else {
				http.Error(w, `{"error": "Ошибка при получении задачи"}`, http.StatusInternalServerError)
			}
			return
		}

		// Если задача не повторяющаяся, удаляем её из базы данных
		if task.Repeat == "" {
			err := storage.RemoveTask(taskID)
			if err != nil {
				http.Error(w, `{"error": "Ошибка при удалении задачи"}`, http.StatusInternalServerError)
				return
			}
		} else {
			// Если задача повторяющаяся, обновляем её следующей датой выполнения
			now := time.Now()
			nextDate, err := nextdate.CalculateNextDate(now, task.Date, task.Repeat)
			if err != nil {
				log.Printf("Ошибка вычисления следующей даты: %v", err)
				http.Error(w, `{"error": "Ошибка вычисления следующей даты"}`, http.StatusInternalServerError)
				return
			}

			// Обновляем задачу в базе данных с новой датой выполнения
			task.Date = nextDate
			err = storage.UpdateTask(task)
			if err != nil {
				http.Error(w, `{"error": "Ошибка при обновлении задачи"}`, http.StatusInternalServerError)
				return
			}
		}

		// Отправляем пустой JSON в ответ на успешное выполнение
		json.NewEncoder(w).Encode(map[string]interface{}{})
	}
}

// UpdateTaskHandler обрабатывает запросы на обновление задачи
func UpdateTaskHandler(storage *database.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var updatedTask task.Task
		if err := json.NewDecoder(r.Body).Decode(&updatedTask); err != nil {
			log.Printf("Ошибка декодирования JSON: %v", err)
			http.Error(w, `{"error": "Некорректный формат данных"}`, http.StatusBadRequest)
			return
		}

		if updatedTask.ID == "" {
			log.Println("Отсутствует идентификатор задачи")
			http.Error(w, `{"error": "Идентификатор задачи не предоставлен"}`, http.StatusBadRequest)
			return
		}

		if err := task.Validate(&updatedTask); err != nil {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
			return
		}

		if err := storage.UpdateTask(updatedTask); err != nil {
			if err.Error() == "задача не найдена" {
				http.Error(w, `{"error": "Задача не найдена"}`, http.StatusNotFound)
			} else {
				http.Error(w, `{"error": "Внутренняя ошибка сервера"}`, http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "success"})
	}
}

// GetTaskHandler обработчик для получения задачи по идентификатору
func GetTaskHandler(storage *database.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		taskID := r.URL.Query().Get("id")
		if taskID == "" {
			http.Error(w, `{"error": "Идентификатор не предоставлен"}`, http.StatusBadRequest)
			return
		}

		task, err := storage.GetTask(taskID)
		if err != nil {
			if err.Error() == "задача не найдена" {
				http.Error(w, `{"error": "Задача не найдена"}`, http.StatusNotFound)
			} else {
				log.Printf("Ошибка при получении задачи: %v", err)
				http.Error(w, `{"error": "Внутренняя ошибка сервера"}`, http.StatusInternalServerError)
			}
			return
		}

		json.NewEncoder(w).Encode(task)
	}
}

// GetTasksHandler обработчик для получения списка задач с возможностью поиска
func GetTasksHandler(storage *database.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		searchTerm := r.URL.Query().Get("search")
		tasks, err := storage.GetTasks(searchTerm, TaskLimit)
		if err != nil {
			log.Printf("Ошибка при получении задач: %v", err)
			http.Error(w, `{"error": "Ошибка выполнения запроса"}`, http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{"tasks": tasks}
		json.NewEncoder(w).Encode(response)
	}
}

// AddTaskHandler обрабатывает запросы на добавление новой задачи
func AddTaskHandler(storage *database.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		var newTask task.Task
		if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
			log.Printf("Ошибка при декодировании JSON: %v", err)
			http.Error(w, `{"error": "Ошибка при декодировании JSON"}`, http.StatusBadRequest)
			return
		}

		// Валидация данных
		if newTask.Title == "" {
			http.Error(w, `{"error": "Title cannot be empty"}`, http.StatusBadRequest)
			return
		}

		taskID, err := storage.InsertTask(newTask)
		if err != nil {
			log.Printf("Ошибка при добавлении задачи в базу данных: %v", err)
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{"id": taskID})
	}
}

// DeleteTaskHandler обработчик для удаления задачи по идентификатору
func DeleteTaskHandler(storage *database.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		taskID := r.URL.Query().Get("id")
		if taskID == "" {
			http.Error(w, `{"error": "Отсутствует идентификатор задачи"}`, http.StatusBadRequest)
			return
		}

		err := storage.RemoveTask(taskID)
		if err != nil {
			if err.Error() == "задача не найдена" {
				http.Error(w, `{"error": "Задача не найдена"}`, http.StatusNotFound)
			} else {
				http.Error(w, `{"error": "Внутренняя ошибка сервера"}`, http.StatusInternalServerError)
			}
			return
		}

		// Возвращаем пустой ответ, так как тесты ожидают пустой результат
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write([]byte(`{}`))
	}
}
