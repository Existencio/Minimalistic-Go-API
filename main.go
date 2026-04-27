package main

import (
	"fmt"
	"net/http"

	"minimalistic-go-api/handlers"
	"minimalistic-go-api/storage"
)

func main() {
	// ! Загружаем посты из файла при старте
	storage.LoadPosts()

	// ! Регистрируем обработчик
	http.HandleFunc("/posts", handlers.PostHandler)
	// ? т.е Когда придёт запрос на /posts, вызови функцию handlers.PostHandler.

	fmt.Println("Server started on port 8080")

	// ! Запускаем сервер
	err := http.ListenAndServe(":8080", nil)
	// ? второй параметр это главный обработчик запросов
	// ? Если передать nil, Go использует стандартный роутер: http.DefaultServeMux
	if err != nil {
		fmt.Println("server error:", err)
	}
}
