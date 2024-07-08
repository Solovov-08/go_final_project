package handlers

import (
	"go_final_project/tasks"
	"net/http"
	"time"
)

// NextDateHandler обрабатывает запросы для получения следующей даты повторения задачи
func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем значения параметров из запроса
	currentDateStr := r.FormValue("now")
	targetDateStr := r.FormValue("date")
	repeatRule := r.FormValue("repeat")

	// Парсим текущее время в формате, указанном в tasks.FormatDate
	currentDate, err := time.Parse(tasks.FormatDate, currentDateStr)
	if err != nil {
		http.Error(w, `{"error": "Неверный формат параметра now"}`, http.StatusBadRequest)
		return
	}

	// Вычисляем следующую дату на основе повторения
	nextDate, err := tasks.CalculateNextDate(currentDate, targetDateStr, repeatRule)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	// Возвращаем следующую дату в виде текста
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write([]byte(`{"next_date": "` + nextDate + `"}`))
}
