package storage

import "url/types"

type Storage interface {
	GetURL(url string) (string, *types.ResponseError)
	CreateShortURL(fullURL, shortURL string) (string, *types.ResponseError)
}
