package storage

import "url/types"

// Интерфейс, определяющий поведение Postgres и InMemory БД
// GetURL - принимает сокращенный URL и возвращает полный URL и ошибку, если та произошла
// CreateShortURL - принимает полный и сокращенный URL, возвращает короткий URL и ошибку, если та произошла
type Storage interface {
	GetURL(shortUrl string) (string, *types.ResponseError)
	IsFullUrlExists(fullUrl string) (bool, string, *types.ResponseError)
	CreateShortURL(fullURL, shortURL string) (string, *types.ResponseError)
}
