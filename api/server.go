package api

import (
	"net/http"
	"url/storage"
)

type Server struct {
	listenAddr string
	storage storage.Storage
}

func NewServer(listenAddr string, storage storage.Storage) *Server {
	return &Server{
		listenAddr: listenAddr,
		storage: storage,
	}
}

func (s *Server) Start() error {
	http.HandleFunc("/", s.requestHandler)
	return http.ListenAndServe(s.listenAddr, nil)
}
