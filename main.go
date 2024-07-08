package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"go_final_project/database"
	"go_final_project/handlers"

	"github.com/go-chi/chi/v5"
)

func main() {
	// Чтение порта из переменной окружения или использование значения по умолчанию
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = strconv.Itoa(7540) // По умолчанию используем порт 7540
	}

	// Инициализация базы данных
	db, err := database.InitializeDatabase()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close() // Закрытие соединения с базой данных при выходе из main()

	// Создание роутера Chi
	r := chi.NewRouter()

	// Обработка статических файлов
	fs := http.FileServer(http.Dir("./web"))
	r.Handle("/*", fs)

	// Установка обработчиков API для операций с задачами
	r.MethodFunc(http.MethodGet, "/api/task", handlers.GetTaskHandler(db))
	r.MethodFunc(http.MethodGet, "/api/tasks", handlers.GetTasksHandler(db))
	r.MethodFunc(http.MethodPut, "/api/task", handlers.UpdateTaskHandler(db))
	r.MethodFunc(http.MethodDelete, "/api/task", handlers.DeleteTaskHandler(db))
	r.MethodFunc(http.MethodPost, "/api/task", handlers.AddTaskHandler(db))
	r.MethodFunc(http.MethodPost, "/api/task/done", handlers.MarkTaskDoneHandler(db))

	// Обработчик API для вычисления следующей даты
	r.HandleFunc("/api/nextdate", handlers.NextDateHandler)

	// Запуск HTTP сервера
	log.Printf("Server is listening on port %s", port)
	err = http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatal(err)
	}
}
