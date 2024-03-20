package api

import (
	"encoding/json"
	"io"
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
		log.Println("Url can't be empty")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Url can't be empty"))
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

	defer r.Body.Close()
	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error while reading body of request"))
		return
	}
	fullURL := string(data)
	shortURL := generateUniqieUrl(letters + fullURL)
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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("http://localhost:8080/" + response))
}
