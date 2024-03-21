package storage

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"url/types"
)

// БД, которая будет храниться в памяти
type MemoryStorage struct {
	memoryDb         map[string]string // Создадим временное хранилище
	memoryDbFileName string            // Имя файла, в котором будут храниться записи
}

// Конструктор БД
func NewMemoryStorage(memoryDbFileName string) *MemoryStorage {
	return &MemoryStorage{
		memoryDb:         map[string]string{},
		memoryDbFileName: memoryDbFileName,
	}
}

// Метод, для проверки, что полный URL уже есть в БД, тогда нам не нужно добавлять новую запись в БД
func (s *MemoryStorage) IsFullUrlExists(fullUrl string) (bool, string, *types.ResponseError) {
	file, err := os.OpenFile(s.memoryDbFileName, os.O_RDONLY, 0666)
	if err != nil {
		log.Println(err.Error())
		return false, "", &types.ResponseError{
			Message: "OpenFile error",
			Status:  http.StatusInternalServerError,
		}
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// Каждую строчку считаем построчно и разделяем по спец. символу - пробелу
		data := strings.Split(scanner.Text(), " ")
		// Пока считываем с файла, будем добавлять во временное хранилище
		s.memoryDb[data[0]] = data[1]
		// Если мы во временное хранилище положили уже больше миллиона записей, то очистим хранилище
		// Чтобы не использовалось слишком много памяти
		if len(s.memoryDb) > 1000000 {
			s.memoryDb = make(map[string]string)
		}

		if data[1] == data[0] {
			return false, data[0], nil
		}
	}
	if scanner.Err() != nil {
		log.Println(scanner.Err())
		return false, "", &types.ResponseError{
			Message: "Scanner error",
			Status:  http.StatusInternalServerError,
		}
	}
	// Если полного URL не найдено, то такой ссылки еще не было в бд
	return true, "", nil
}

// Метод для получения полного URL
// В параметры метода передается короткий URL
func (s *MemoryStorage) GetURL(url string) (string, *types.ResponseError) {
	// Проверим, был ли когда выполнен к данному URL запрос
	if fullURL, ok := s.memoryDb[url]; !ok {
		// Если URL нету в map, то будем читать построчно с файла
		file, err := os.OpenFile(s.memoryDbFileName, os.O_RDONLY, 0666)
		if err != nil {
			log.Println(err.Error())
			return "", &types.ResponseError{
				Message: "OpenFile error",
				Status:  http.StatusInternalServerError,
			}
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			// Каждую строчку считаем построчно и разделяем по спец. символу - пробелу
			data := strings.Split(scanner.Text(), " ")
			// Если мы во временное хранилище положили уже больше миллиона записей, то очистим хранилище
			// Чтобы не использовалось слишком много памяти
			if len(s.memoryDb) > 1000000 {
				s.memoryDb = make(map[string]string)
			}
			// Если короткий URL есть в файле, то вернем его полный URL и добавим в map
			if data[0] == url {
				s.memoryDb[data[0]] = data[1]
				return data[1], nil
			}
		}
		if scanner.Err() != nil {
			log.Println(scanner.Err())
			return "", &types.ResponseError{
				Message: "Scanner error",
				Status:  http.StatusInternalServerError,
			}
		}
		// Если короткого URL не найдено, то вернем ошибку
		return "", &types.ResponseError{
			Message: "Invalid Short url. Not found.",
			Status:  http.StatusBadRequest,
		}
	} else {
		// В map лежит короткий URL, тогда просто вернем полный URL
		return fullURL, nil
	}
}

// Метод для создания короткого URL
// В параметры метода передадим полный и короткий URL сайта
func (s *MemoryStorage) CreateShortURL(fullURL, shortURL string) (string, *types.ResponseError) {
	// Откроем файл на чтение и запись
	file, err := os.OpenFile(s.memoryDbFileName, os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		log.Println("Error while opening inMemoryDb")
		return "", &types.ResponseError{
			Message: "Error while opening inMemoryDb",
			Status:  http.StatusInternalServerError,
		}
	}
	defer file.Close()

	if _, err := fmt.Fprintln(file, shortURL+" "+fullURL); err != nil {
		log.Println("Error while appending new line in dbMemory")
		return "", &types.ResponseError{
			Message: "Error while appending new line in dbMemory",
			Status:  http.StatusInternalServerError,
		}
	}
	// Добавим во временное хранилище запись и вернем хендлеру
	s.memoryDb[shortURL] = fullURL
	// Если мы во временное хранилище положили уже больше миллиона записей, то очистим хранилище
	// Чтобы не использовалось слишком много памяти
	if len(s.memoryDb) > 1000000 {
		s.memoryDb = make(map[string]string)
	}
	return shortURL, nil
}
