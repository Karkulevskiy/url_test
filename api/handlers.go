package api

import (
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"net/http"
)

// Перемешанные доступные символы для создания уникального короткого URL
const setOfChars = "Ot43qYsT6ZzDABXh9S05g8PdeMwJV71lumkHFnQKboLCafUGWcN82R4IvixyjpEGr"

// HTTP Handler, который определяет был сделан POST или GET запрос
func (s *Server) requestHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.GetURL(w, r) // Был сделан GET запрос
	case http.MethodPost:
		s.CreateShortURL(w, r) // Был сделан POST запрос
	}
}

// Функция для генерации рандомного URL
func (s *Server) generateUniqueURL(template string) string {
	var uniqueUrl []byte
	// Будет создан URL из 6 символов
	// Если такой короткий URL уже суещствует, то будем добавлять ему еще по 6
	// рандомных символов, пока не создадим уникальный
	for {
		for i := 0; i < 6; i++ {
			uniqueUrl = append(uniqueUrl, template[rand.Intn(len(template))])
		}
		resp, _ := s.storage.GetURL(string(uniqueUrl))
		if resp == "" {
			return string(uniqueUrl)
		}
	}
}

// HTTP GET Handler
func (s *Server) GetURL(w http.ResponseWriter, r *http.Request) {
	// Мы получаем URL сайта через строку запроса, поэтому просто возьмем ее от туда
	url := r.URL.Path[1:]
	// Проверим на пустую строку
	if url == "" {
		log.Println("Url can't be empty")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Url can't be empty"))
		return
	}
	// Получаем полный URL с БД
	response, err := s.storage.GetURL(url)
	if err != nil {
		// Если произошла ошибка, то вернем ответ об ошибке
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(err.Status)
		jsonResp, err := json.Marshal(err)
		if err != nil {
			log.Println("Error while marshaling json")
			return
		}
		w.Write(jsonResp)
		return
	}
	// Устанавливаем Header, body для ответа
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}

// HTTP Handler для создания короткого короткого URL
func (s *Server) CreateShortURL(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	// Считываем с тела запроса введенный сайт
	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error while reading body of request"))
		return
	}
	// Проверка, чтобы тело запроса не было пустым
	if data == nil {
		log.Println("Error, body can't be empty")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error, body can't be empty"))
		return
	}
	fullURL := string(data)
	// Создаем уникальную URL
	shortURL := s.generateUniqueURL(setOfChars)
	// Добавим в БД новую сущность с полным и коротким URL
	// Проверим на ошибку
	response, er := s.storage.CreateShortURL(fullURL, shortURL)
	if er != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(er.Status)
		jsonResp, er := json.Marshal(er)
		if er != nil {
			log.Println("Error while marshaling json")
			return
		}
		w.Write(jsonResp)
		return
	}
	// Отправим ответ с соответсвующем Header, Body
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("http://localhost:8080/" + response))
}
