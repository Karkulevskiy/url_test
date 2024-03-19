package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	_ "github.com/lib/pq"
)

var letters = "Ot43qYsT6ZzDABXh9S05g8PdeMwJV71lumkHFnQKboLCafUGWcN82R4IvixyjpEGr"

func generateUniqieUrl(template string) string {
	var uniqueUrl []rune
	for i := 0; i < 6; i++ {
		uniqueUrl = append(uniqueUrl, rune(template[rand.Intn(len(template))]))
	}
	return string(uniqueUrl)
}

func GetUrl(w http.ResponseWriter, r *http.Request) {
	fmt.Println("get")
	url := r.URL.Path[1:]
	connStr := "host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Println("db error")
		fmt.Println(err.Error())
		return
	}
	defer db.Close()
	query := `select full_url from urls where short_url = $1`
	fmt.Println(url)
	rows, err := db.Query(query, url)
	if err != nil {
		log.Println("error while makeing query")
		log.Println(err.Error())
		return
	}
	defer rows.Close()
	var fullUrl string
	for rows.Next() {
		err := rows.Scan(&fullUrl)
		if err != nil {
			log.Println("not found")
			log.Println(err.Error())
			return
		}
	}
	if rows.Err() != nil {
		log.Println("last err")
		return
	}
	//make response
	fmt.Println(fullUrl)
}

func CutURL(w http.ResponseWriter, r *http.Request) {
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
	//fmt.Println(res.RowsAffected())
}

func requestHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetUrl(w, r)
	case http.MethodPost:
		CutURL(w, r)
	}
}

func main() {
	//http.HandleFunc("/", CutURL)
	//http.HandleFunc("/", GetUrl)
	http.HandleFunc("/", requestHandler)
	http.ListenAndServe("localhost:8080", nil)
}
