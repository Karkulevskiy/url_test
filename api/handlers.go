package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
)

const letters = "Ot43qYsT6ZzDABXh9S05g8PdeMwJV71lumkHFnQKboLCafUGWcN82R4IvixyjpEGr"

func (s *Server) requestHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.GetURL(w, r)
	case http.MethodPost:
		s.CreateShortURL(w, r)
	}
}

func generateUniqieUrl(template string) string {
	var uniqueUrl []rune
	for i := 0; i < 6; i++ {
		uniqueUrl = append(uniqueUrl, rune(template[rand.Intn(len(template))]))
	}
	return string(uniqueUrl)
}

func (s *Server) GetURL(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path[1:]
	if url == "" {
		//обработать данный момент с null == url
		log.Println("Url can't be empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	response, err := s.storage.GetURL(url)
	if err != nil {
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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}

func (s *Server) CreateShortURL(w http.ResponseWriter, r *http.Request) {
	fmt.Println("post")
	url := r.URL.Path[1:]
	connStr := "host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Println("db error")
		fmt.Println(err.Error())
		return
	}
	defer db.Close()

	template := letters + url

	uniqueUrl := generateUniqieUrl(template)

	res, err := db.Exec("insert into urls (full_url, short_url) values ($1, $2)", url, uniqueUrl)
	if err != nil {
		log.Println("insertin err")
		fmt.Println(err.Error())
		return
	}
	_ = res
}
