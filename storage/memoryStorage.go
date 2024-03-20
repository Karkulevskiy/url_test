package storage

import (
	"net/http"
	"url/types"
)

type MemoryStorage struct {
	memoryDb map[string]string
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		memoryDb: map[string]string{},
	}
}

func (s *MemoryStorage) GetURL(url string) (string, *types.ResponseError) {
	if fullURL, ok := s.memoryDb[url]; !ok {
		return "", &types.ResponseError{
			Message: "Invalid Short url. Not found.",
			Status:  http.StatusBadRequest,
		}
	} else {
		return fullURL, nil
	}
}

func (s *MemoryStorage) CreateShortURL(fullURL, shortURL string) (string, *types.ResponseError) {
	s.memoryDb[shortURL] = fullURL
	return shortURL, nil
}
