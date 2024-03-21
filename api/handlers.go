package api

import (
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"net/http"
)

const (
	setBytes      = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789" // Набор всех возможных символов, при генерации URL
	letterIdxBits = 6                                                                // 6 битов для создания маски
	letterIdxMask = 1<<letterIdxBits - 1                                             // Маска для побитового И
	letterIdxMax  = 63 / letterIdxBits                                               // Будем хранить остаток битов, чтобы снова не генерировать рандомное число
)

// HTTP Handler, который определяет был сделан POST или GET запрос
func (s *Server) requestHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.GetURL(w, r) // Был сделан GET запрос
	case http.MethodPost:
		s.CreateShortURL(w, r) // Был сделан POST запрос
	}
}

// Функция для генерации рандомного короткого URL
// Параметром функции передается длина короткого URL
func (s *Server) generateUniqueURL(n int) string {
	b := make([]byte, n)
	// i - счетчик сгенерированных симолов
	// cache - переменная, где будем хранить оставшиеся неиспользованные биты
	// remain - сколько осталось неиспользованных битов
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		// Если все уже были использованы, то обновим их
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		// Будем делать побитовое AND, если idx больше длины setBytes, то наложим снова маску
		if idx := int(cache & letterIdxMask); idx < len(setBytes) {
			b[i] = setBytes[idx] // Добавим в короткие URL символ
			i--
		}
		cache >>= letterIdxBits // Побитовый сдвиг, чтобы взять оставшиеся неиспользованные биты
		remain--                // Уменьшим счетчик оставшихся битов
	}

	return string(b)
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

	// Проверим, существует ли в БД уже такой полный URL
	isUnique, shortUrl, err_ := s.storage.IsFullUrlExists(fullURL)

	// Проверка, если произошла ошибка во время запроса к БД
	if err_ != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(err_.Status)
		jsonResp, er := json.Marshal(err_)
		if er != nil {
			log.Println("Error while marshaling json")
			return
		}
		w.Write(jsonResp)
		return
	}

	// Если уже полная URL, то вернем короткую ссылку
	if !isUnique {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("http://localhost:8080/" + shortUrl))
		return
	}
	//Если полного URL в БД не нашлось, то создадимя для него короткую ссылку
	shortURL := s.generateUniqueURL(6)
	// Проверка, что такого короткого уникального URL еще не создавали
	resp, _ := s.storage.GetURL(shortURL)
	// Если оказывается, что такой URL уже существует
	// То будем к короткой ссылке добавлять еще по 3 символа
	// Вероятность создать два одинаковых URL +- == 1,7605561 × 10^−11
	if resp == shortURL {
		for {
			shortURL = shortURL + s.generateUniqueURL(3)
			resp, _ = s.storage.GetURL(shortURL)
			if resp != shortURL {
				break
			}
		}
	}
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
