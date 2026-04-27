package handlers

import (
	"encoding/json"
	"net/http"

	"minimalistic-go-api/models"
	"minimalistic-go-api/storage"
)

// ! Функция-обработчик HTTP запросов
func PostHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet: // ? http.MethodGet это готовая константа в го, которая всегда равна строке "GET"
		posts := storage.GetPosts() // ! создаём ещё одну переменную posts уже локальную для этого обработчика
		json.NewEncoder(w).Encode(posts)

	case http.MethodPost:
		var newPost models.Post
		// возьми JSON из тела запроса и запиши данные внутрь newPost
		err := json.NewDecoder(r.Body).Decode(&newPost) // ? r.Body это тело HTTP запроса
		if err != nil {
			w.WriteHeader(http.StatusBadRequest) // ! http.StatusBadRequest это константа в го, которая равна 400
			// ? отправить клиенту JSON с описанием ошибки
			json.NewEncoder(w).Encode(map[string]string{
				"error": "invalid JSON",
			})
			return
		}

		err = storage.AddPost(newPost)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "failed to save post",
			})

			return
		}

		w.WriteHeader(http.StatusCreated) // ! http.StatusCreated константа = 201 (запрос успешен, новый ресурс создан)
		json.NewEncoder(w).Encode(map[string]any{
			"message": "post created",
			"data":    newPost,
		})

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "method not allowed",
		})
	}
}
