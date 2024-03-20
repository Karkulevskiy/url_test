package main

import (
	"fmt"
	"log"
	"url/api"
	"url/storage"
	_ "github.com/lib/pq"
)

func main() {
	listenAddr := "localhost:8080"
	storage := storage.Postgres{}
	server := api.NewServer(listenAddr, &storage)
	fmt.Println("Server is running on port:", listenAddr)
	log.Fatal(server.Start())
}
