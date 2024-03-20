package main

import (
	"fmt"
	"log"
	"os"
	"url/api"
	"url/storage"

	_ "github.com/lib/pq"
)

const listenAddr = "localhost:8080"

func main() {
	var server *api.Server
	var db storage.Storage
	switch os.Args[1] {
	case "-d":
		db = storage.NewPostgres()
	default:
		db = storage.NewMemoryStorage()
	}
	server = api.NewServer(listenAddr, db)
	fmt.Println("Server is running on port:", listenAddr)
	log.Fatal(server.Start())
}
