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
	connStr              string // Строка подключения к БД
	queryGet             string // SELECT запрос к БД для получения адреса сайта
	queryInsert          string // INSERT запрос к БД для добавления адреса сайта
	queryIsFullUrlUnique string // SELECT запрос к БД на проверку уникальности полного URL
}

// Констуктор для Postgres
// ConnStr - строка подключения к БД
func NewPostgres(connStr string) *Postgres {
	return &Postgres{
		connStr:              connStr,
		queryGet:             "SELECT full_url FROM urls WHERE short_url = $1",
		queryInsert:          "INSERT INTO urls (full_url, short_url) VALUES ($1, $2) RETURNING short_url",
		queryIsFullUrlUnique: "SELECT short_url FROM urls WHERE full_url = $1",
	}
}

// Метод, для проверки, что полный URL уже есть в БД, тогда нам не нужно добавлять новую запись в БД
func (s *Postgres) IsFullUrlExists(fullUrl string) (bool, string, *types.ResponseError) {
	// Подключение к БД
	db, err := sql.Open("postgres", s.connStr)
	if err != nil {
		log.Println(err.Error())
		return false, "", &types.ResponseError{
			Message: "Can't open DB",
			Status: http.StatusInternalServerError,
		}
	}
	defer db.Close()
	// Выполнения запроса к БД
	rows, err := db.Query(s.queryIsFullUrlUnique, fullUrl)
	if err != nil {
		log.Println(err.Error())
		return false, "", &types.ResponseError{
			Message: "Querry error",
			Status: http.StatusInternalServerError,
		}
	}
	defer rows.Close()
	var URL string
	// Получение полного адреса сайта
	for rows.Next() {
		err = rows.Scan(&URL)
		if err != nil {
			log.Println(err.Error())
			return false, "", &types.ResponseError{
				Message: "Scan error",
				Status: http.StatusInternalServerError,
			}
		}
	}
	if rows.Err() != nil {
		log.Println(rows.Err().Error())
		return false, "", &types.ResponseError{
			Message: "rows Error",
			Status: http.StatusInternalServerError,
		}
	}
	// Проверка, если сайт по полному URL не найден
	if URL == "" {
		return true, "", nil
	}
	return false, URL, nil	  
}


// Метод, для получения полного или короткого URL из БД
func (s *Postgres) GetURL(shortUrl string) (string, *types.ResponseError) {
	// Подключение к БД
	db, err := sql.Open("postgres", s.connStr)
	if err != nil {
		log.Println(err.Error())
		return "", &types.ResponseError{
			Message: "Can't connect to postgres",
			Status:  http.StatusInternalServerError,
		}
	}
	defer db.Close()
	// Выполнения запроса к БД
	rows, err := db.Query(s.queryGet, shortUrl)
	if err != nil {
		log.Println(err.Error())
		return "", &types.ResponseError{
			Message: "Query error",
			Status:  http.StatusInternalServerError,
		}
	}
	defer rows.Close()
	var URL string
	// Получение полного адреса сайта
	for rows.Next() {
		err = rows.Scan(&URL)
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
	// Проверка, если сайт по короткому URL не найден
	if URL == "" {
		return "", &types.ResponseError{
			Message: "Invalid Short url. Not Found",
			Status:  http.StatusBadRequest,
		}
	}
	return URL, nil
}

// Метод, для добавления сайта к БД
// Аргументы метода: полная и короткая ссылка
func (s *Postgres) CreateShortURL(fullURL, shortURL string) (string, *types.ResponseError) {
	// Подключения к БД
	db, err := sql.Open("postgres", s.connStr)
	if err != nil {
		log.Println(err.Error())
		return "", &types.ResponseError{
			Message: "Can't connect to postgres",
			Status:  http.StatusInternalServerError,
		}
	}
	defer db.Close()
	// Выполнение запроса
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
	// Получение короткой ссылки
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
