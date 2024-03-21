package main

import (
	"fmt"
	"log"
	"os"
	"url/api"
	"url/storage"

	_ "github.com/lib/pq"
)

const (
	// TCP network address
	listenAddr = "localhost:8080"
	// Строка подлючения для Postgres
	connStr          = "host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable"
	memoryDbFileName = "memoryDb.txt"
)

func main() {
	var server *api.Server
	var db storage.Storage
	if len(os.Args) > 1 && os.Args[1] == "-d" {
		db = storage.NewPostgres(connStr) // Получаем экземпляр Postgres
	} else {
		pwd, err := os.Getwd()
		if err != nil {
			log.Println("Can't get working directory")
			log.Println(err.Error())
			return
		}
		path := pwd + "/" + memoryDbFileName
		file, err := os.OpenFile(path, os.O_CREATE, 0666)
		if err != nil{
			log.Println("Error while creating memory DB")
			log.Println(err.Error())
			return
		}
		file.Close() 
		db = storage.NewMemoryStorage(path) // Создаем БД в памяти
	}
	server = api.NewServer(listenAddr, db) // Получаем экземпляр сервера
	fmt.Println("Server is running on port:", listenAddr)
	log.Fatal(server.Start()) // Запускаем сервер
}
