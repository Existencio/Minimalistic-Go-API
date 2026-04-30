package storage

import (
	"encoding/json"
	"os"

	"minimalistic-go-api/models"
)

// ! Слайс, в котором будут храниться посты
var posts []models.Post

// ? Скобки [] указывают на тип данных slice
// ? models - пакет, откуда берём тип
// ? Post - сам тип. Наша пользовательская структура из пакета models

// ! Задаём унифицированное название файла
const fileName = "posts.json"

// ! Функция, загружающая посты из файла posts.json
// ! С БОЛЬШОЙ БУКВЫ Т.К. ФУНКЦИЯ ЭКСПОРТИРУЕТСЯ
func LoadPosts() {
	data, err := os.ReadFile(fileName) // ? os.ReadFile возвращает data типа []byte и err типа error
	if err != nil {                    // ? err хранит либо nil при удаче, либо значение error при ошибке
		posts = []models.Post{} // ? Если файл прочитать не получилось, делаем posts пустым списком постов
		return                  // ? Выйти из функции прямо сейчас
	}

	// ? Если же файл прочитался
	// ! Превращаем JSON из файла в Go структуры
	err = json.Unmarshal(data, &posts)
	if err != nil { // ? Если JSON не получилось превратить в posts
		posts = []models.Post{} // ? Делаем posts пустым списком постов
		return                  // ? Выйти из функции прямо сейчас
	}
	// ? Используем именно указатель &posts
	// ? т.к. json.Unmarshal должен не просто прочитать данные,
	// ! а записать результат ВНУТРЬ переменной posts
	// ? для этого ему нужен доступ к самой переменной, а не её копии
}

// ! Функция, сохраняющая текущий список posts в файл posts.json
// ? т.е. её можно вызвать из других пакетов, например из handlers
func SavePosts() error { // ? error - тип возвращаемого значения
	// ? json.MarshalIndent делает из Go-данных красивый json с отступами
	// ? data имеет тип []byte, то есть срез байт
	// ? err хранит либо nil при удаче, либо значение error при ошибке
	data, err := json.MarshalIndent(posts, "", "  ")
	if err != nil { // ? Если не получилось превратить posts в JSON
		// ! ВАЖНО: здесь мы НЕ обрабатываем ошибку
		// ? мы только отдаём её тому месту, где вызвали SavePosts()
		return err // ? Возвращаем ошибку наружу
	}

	// ! Записываем data в файл posts.json
	// ? 0644 это права доступа к файлу
	// ? владелец может читать и писать, остальные читать
	// ! Если os.WriteFile вернёт ошибку, SavePosts отдаст её наружу
	// ? обрабатывать эту ошибку нужно там, где вызвали SavePosts()
	return os.WriteFile(fileName, data, 0644)
}

// ! Функция, отдающая текущий список постов из памяти
// ? а именно из переменной posts
func GetPosts() []models.Post { // ? []models.Post это тип возвращаемого значения
	return posts
}

// ! Функция, добавляющая новый пост в список постов
// ! и пытающаяся сохранить список в файл
func AddPost(post models.Post) error { // ? post - имя параметра, models.Post - тип параметра
	// ? append(куда_добавить, что_добавить)
	// ? не меняет текущий слайс напрямую, но возвращает НОВЫЙ слайс
	// ! поэтому его надо переприсвоить
	posts = append(posts, post)

	// ! Сохраняем обновлённый список posts в файл
	// ? SavePosts возвращает nil при удаче или error при ошибке
	// ! AddPost тоже отдаёт эту ошибку наружу
	// ? обрабатывать её нужно там, где вызвали AddPost()
	return SavePosts()
}

func GetPostsByID(id int) (models.Post, bool) { // ? 1 скобки - параметр, 2 - возвращаемые значения
	for _, post := range posts { // ? при range Go даёт два значения index, value
		if post.ID == id {
			return post, true
		}
	}
	return models.Post{}, false // ? пустой объект
}

func UpdatePost(id int, updatedPost models.Post) (models.Post, bool, error) {
	for i, post := range posts {
		if post.ID == id {
			updatedPost.ID = id
			posts[i] = updatedPost

			err := SavePosts()
			if err != nil {
				return models.Post{}, false, err
			}

			return updatedPost, true, nil
		}
	}

	return models.Post{}, false, nil
}

func DeletePost(id int) (bool, error) {
	for i, post := range posts {
		if id == post.ID {
			posts = append(posts[:i], posts[i+1:]...)
			// ? ... распаковка слайса позволяет передать посты как отдельные элементы
			// ? потому что для элементов слайса задан тип Post, а без ... будут элементы с типом слайс

			err := SavePosts()
			if err != nil {
				return false, err
			}

			return true, nil
		}
	}

	return false, nil
}
