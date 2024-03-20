package api

import (
	"net/http"
	"url/storage"
)
//Сервер, хранящий адрес прослушивания и БД
type Server struct {
	listenAddr string
	storage storage.Storage
}

//Конструктор для создания сервера
// Конструктор принимает адрес прослушивания и БД (Postgres или в InMemory)
func NewServer(listenAddr string, storage storage.Storage) *Server {
	return &Server{
		listenAddr: listenAddr,
		storage: storage,
	}
}

//Метод для запуска сервера
func (s *Server) Start() error {
	http.HandleFunc("/", s.requestHandler)
	return http.ListenAndServe(s.listenAddr, nil)
}
