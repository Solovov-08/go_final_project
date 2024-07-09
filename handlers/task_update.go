package handlers

import (
	"encoding/json"
	"go_final_project/database"
	"go_final_project/task"
	"log"
	"net/http"
)

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
