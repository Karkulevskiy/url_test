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
	//TCP network address
	listenAddr = "localhost:8080"
	//Строка подлючения для Postgres
	connStr = "host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable"
)

func main() {
	var server *api.Server
	var db storage.Storage
	switch os.Args[1] { //Параметры запуска
	case "-d": 
		db = storage.NewPostgres(connStr) //Получаем экземпляр Postgres
	default:
		db = storage.NewMemoryStorage() //Создаем БД в памяти
	}
	server = api.NewServer(listenAddr, db) //Получаем экземпляр сервера
	fmt.Println("Server is running on port:", listenAddr)
	log.Fatal(server.Start()) //Запускаем сервер
}
