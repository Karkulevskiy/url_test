package storage

import (
	"database/sql"
	"log"
	"net/http"
	"url/types"

	_ "github.com/lib/pq"
)

type Postgres struct {
	connStr     string
	queryGet    string
	queryInsert string
}

func NewPostgres() *Postgres {
	return &Postgres{
		connStr:     "host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable",
		queryGet:    "select full_url from urls where short_url = $1",
		queryInsert: "insert into urls (full_url, short_url) values ($1, $2) returning short_url",
	}
}

func (s *Postgres) GetURL(url string) (string, *types.ResponseError) {
	db, err := sql.Open("postgres", s.connStr)
	if err != nil {
		log.Println(err.Error())
		return "", &types.ResponseError{
			Message: "Can't connect to postgres",
			Status:  http.StatusInternalServerError,
		}
	}
	defer db.Close()
	rows, err := db.Query(s.queryGet, url)
	if err != nil {
		log.Println(err.Error())
		return "", &types.ResponseError{
			Message: "Query error",
			Status:  http.StatusInternalServerError,
		}
	}
	defer rows.Close()
	var fullURL string
	for rows.Next() {
		err = rows.Scan(&fullURL)
		if err != nil {
			log.Println(err.Error())
			return "", &types.ResponseError{
				Message: "rows.Scan() error",
				Status:  http.StatusInternalServerError,
			}
		}
	}
	if rows.Err() != nil {
		log.Println(rows.Err().Error())
		return "", &types.ResponseError{
			Message: "rows.Scan() error",
			Status:  http.StatusInternalServerError,
		}
	}
	if fullURL == "" {
		return "", &types.ResponseError{
			Message: "Invalid Short url. Not Found",
			Status:  http.StatusBadRequest,
		}
	}
	return fullURL, nil
}

func (s *Postgres) CreateShortURL(fullURL, shortURL string) (string, *types.ResponseError) {
	db, err := sql.Open("postgres", s.connStr)
	if err != nil {
		log.Println(err.Error())
		return "", &types.ResponseError{
			Message: "Can't connect to postgres",
			Status:  http.StatusInternalServerError,
		}
	}
	defer db.Close()
	rows, err := db.Query(s.queryInsert, fullURL, shortURL)
	if err != nil {
		log.Println(rows.Err())
		return "", &types.ResponseError{
			Message: "Error while inserting into db",
			Status:  http.StatusInternalServerError,
		}
	}
	defer rows.Close()
	var shortURLResponse string
	for rows.Next() {
		err = rows.Scan(&shortURLResponse)
		if err != nil {
			log.Println(err.Error())
			return "", &types.ResponseError{
				Message: "Error while rows.Scan() shortURLResponse",
				Status:  http.StatusInternalServerError,
			}
		}
	}
	if rows.Err() != nil {
		log.Println(rows.Err().Error())
		return "", &types.ResponseError{
			Message: "Error while rows.Scan() shortURLResponse",
			Status:  http.StatusInternalServerError,
		}
	}
	return shortURLResponse, nil
}
