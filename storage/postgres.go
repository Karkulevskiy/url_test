package storage

import (
	"database/sql"
	"log"
	"net/http"
	"url/types"

	_ "github.com/lib/pq"
)

// Структура для описания БД Postgres
type Postgres struct {
	connStr     string //Строка подключения к БД
	queryGet    string //SELECT запрос к БД для получения адреса сайта
	queryInsert string //INSERT запрос к БД для добавления адреса сайта
}

// Констуктор для Postgres
// ConnStr - строка подключения к БД
func NewPostgres(connStr string) *Postgres {
	return &Postgres{
		connStr:     connStr,
		queryGet:    "select full_url from urls where short_url = $1",
		queryInsert: "insert into urls (full_url, short_url) values ($1, $2) returning short_url",
	}
}

// Метод, принимающий через url query короткую ссылку
// для получения в БД полной ссылки сайта
func (s *Postgres) GetURL(url string) (string, *types.ResponseError) {
	//Подключение к БД
	db, err := sql.Open("postgres", s.connStr)
	if err != nil {
		log.Println(err.Error())
		return "", &types.ResponseError{
			Message: "Can't connect to postgres",
			Status:  http.StatusInternalServerError,
		}
	}
	defer db.Close()
	//Выполнения запроса к БД
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
	//Получение полного адреса сайта
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
	//Проверка, если сайт по короткому URL не найден
	if fullURL == "" {
		return "", &types.ResponseError{
			Message: "Invalid Short url. Not Found",
			Status:  http.StatusBadRequest,
		}
	}
	return fullURL, nil
}

// Метод, для добавления сайта к БД
// Аргументы метода: полная и короткая ссылка
func (s *Postgres) CreateShortURL(fullURL, shortURL string) (string, *types.ResponseError) {
	//Подключения к БД
	db, err := sql.Open("postgres", s.connStr)
	if err != nil {
		log.Println(err.Error())
		return "", &types.ResponseError{
			Message: "Can't connect to postgres",
			Status:  http.StatusInternalServerError,
		}
	}
	defer db.Close()
	//Выполнение запроса
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
	//Получение короткой ссылки
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
