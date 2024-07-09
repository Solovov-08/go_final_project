package handlers

import (
	"encoding/json"
	"go_final_project/database"
	"go_final_project/task"
	"log"
	"net/http"
)

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
