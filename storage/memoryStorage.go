package storage

import (
	"net/http"
	"url/types"
)
//БД, которая будет храниться в памяти
type MemoryStorage struct {
	memoryDb map[string]string //Создадим map
}

//Конструктор БД
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		memoryDb: map[string]string{},
	}
}

//Метод для получения полного URL
// В параметры метода передается короткий URL
func (s *MemoryStorage) GetURL(url string) (string, *types.ResponseError) {
	if fullURL, ok := s.memoryDb[url]; !ok {
		//Если короткого URL не найдено, то вернем ошибку
		return "", &types.ResponseError{
			Message: "Invalid Short url. Not found.",
			Status:  http.StatusBadRequest,
		}
	} else {
		return fullURL, nil
	}
}

//Метод для создания короткого URL
// В параметры метода передадим полный и короткий URL сайта
func (s *MemoryStorage) CreateShortURL(fullURL, shortURL string) (string, *types.ResponseError) {
	s.memoryDb[shortURL] = fullURL
	return shortURL, nil
}
