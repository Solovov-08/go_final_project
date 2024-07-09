package handlers

import (
	"encoding/json"
	"go_final_project/database"
	"log"
	"net/http"
)

const TaskLimit = 50

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
