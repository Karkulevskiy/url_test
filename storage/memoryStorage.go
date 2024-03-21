package storage

import (
	"bufio"
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

// Метод для получения полного URL
// В параметры метода передается короткий URL
func (s *MemoryStorage) GetURL(url string) (string, *types.ResponseError) {
	// Проверим, был ли когда выполнен к данному URL запрос
	if fullURL, ok := s.memoryDb[url]; !ok {
		// Если URL нету в map, то будем читать построчно с файла
		file, err := os.OpenFile(s.memoryDbFileName, os.O_RDONLY, 0666)
		if err != nil {
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
			log.Println("Error while scanning file")
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
	file, err := os.OpenFile(s.memoryDbFileName, os.O_RDWR, 0666)
	if err != nil {
		log.Println("Error while opening inMemoryDb")
		return "", &types.ResponseError{
			Message: "Error while opening inMemoryDb",
			Status:  http.StatusInternalServerError,
		}
	}
	defer file.Close()
	// Будем искать уже существующий полный URL, чтобы не делать перезапись
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data := strings.Split(scanner.Text(), " ")
		//Сравниваем полный URL из файла с полным URL аргументом метода
		if data[1] == fullURL {
			//Вернем уже сущесвтующий shortURL и добавим во временное хранилище
			s.memoryDb[data[0]] = data[1]
			return data[0], nil
		}
	}
	// Проверим на возникновение ошибок
	if scanner.Err() != nil {
		log.Println("Error while searching existed URL")
		return "", &types.ResponseError{
			Message: "Error while searching existed URL",
			Status:  http.StatusInternalServerError,
		}
	}
	// Если мы не вышли из метода, то такой записи еще не было
	// Тогда запишем в файл shortURL, fullURL
	if _, err := file.WriteString(shortURL + " " + fullURL); err != nil {
		log.Println("Error while appending new line in dbMemory")
		return "", &types.ResponseError{
			Message: "Error while appending new line in dbMemory",
			Status:  http.StatusInternalServerError,
		}
	}
	// Если мы во временное хранилище положили уже больше миллиона записей, то очистим хранилище
	// Чтобы не использовалось слишком много памяти
	if len(s.memoryDb) > 1000000 {
		s.memoryDb = make(map[string]string)
	}
	// Добавим во временное хранилище запись и вернем хендлеру
	s.memoryDb[shortURL] = fullURL
	return shortURL, nil
}
