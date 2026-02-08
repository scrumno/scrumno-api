package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	// 1. Описываем маршрут (роут).
	// При обращении к корню "/" будет срабатывать эта функция
	http.HandleFunc("/health/check", func(w http.ResponseWriter, r *http.Request) {
		// Устанавливаем заголовок, что мы отправляем JSON
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")

		// Создаем данные для ответа
		response := map[string]string{"message": "hello world"}

		// Кодируем данные в JSON и отправляем клиенту
		json.NewEncoder(w).Encode(response)
	})

	log.Println("Сервер запущен на порту :8080")

	// 2. Запускаем сервер на порту 8080
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}