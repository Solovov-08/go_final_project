package handlers

import (
	"encoding/json"
	"go_final_project/database"
	"go_final_project/model"
	"go_final_project/tasks"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
)

// UpdateTaskHandler обрабатывает запросы на обновление задачи
func UpdateTaskHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var updatedTask model.Task
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

		if err := tasks.ValidateTask(&updatedTask); err != nil {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
			return
		}

		if err := database.UpdateTask(db, updatedTask); err != nil {
			if err.Error() == "task not found" {
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
