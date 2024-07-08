package handlers

import (
	"encoding/json"
	"go_final_project/database"
	"go_final_project/model"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
)

// AddTaskHandler обрабатывает запросы на добавление новой задачи
func AddTaskHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		var newTask model.Task
		if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
			log.Printf("Ошибка при декодировании JSON: %v", err)
			http.Error(w, `{"error": "Ошибка при декодировании JSON"}`, http.StatusBadRequest)
			return
		}

		taskID, err := database.InsertTask(db, newTask)
		if err != nil {
			log.Printf("Ошибка при добавлении задачи в базу данных: %v", err)
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{"id": taskID})
	}
}

// DeleteTaskHandler обработчик для удаления задачи по идентификатору
func DeleteTaskHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		taskID := r.URL.Query().Get("id")
		if taskID == "" {
			http.Error(w, `{"error": "Отсутствует идентификатор задачи"}`, http.StatusBadRequest)
			return
		}

		err := database.RemoveTask(db, taskID)
		if err != nil {
			if err.Error() == "task not found" {
				http.Error(w, `{"error": "Задача не найдена"}`, http.StatusNotFound)
			} else {
				http.Error(w, `{"error": "Внутренняя ошибка сервера"}`, http.StatusInternalServerError)
			}
			return
		}

		//json.NewEncoder(w).Encode(map[string]interface{}{"status": "success"})

		// Возвращаем пустой ответ, так как тесты ожидают пустой результат
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write([]byte(`{}`))
	}
}
