package handlers

import (
	"encoding/json"
	"go_final_project/database"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
)

const TaskLimit = 50

// GetTaskHandler обработчик для получения задачи по идентификатору
func GetTaskHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		taskID := r.URL.Query().Get("id")
		if taskID == "" {
			http.Error(w, `{"error": "Идентификатор не предоставлен"}`, http.StatusBadRequest)
			return
		}

		task, err := database.FetchTaskByID(db, taskID)
		if err != nil {
			if err.Error() == "task not found" {
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
func GetTasksHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		searchTerm := r.URL.Query().Get("search")
		tasks, err := database.FetchTasks(db, searchTerm, TaskLimit)
		if err != nil {
			log.Printf("Ошибка при получении задач: %v", err)
			http.Error(w, `{"error": "Ошибка выполнения запроса"}`, http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{"tasks": tasks}
		json.NewEncoder(w).Encode(response)
	}
}
