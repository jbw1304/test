package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/mattn/go-sqlite3"
)
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func shorting() string {
	b := make([]byte, 5)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
func isValidUrl(token string) bool {
	_, err := url.ParseRequestURI(token)
	if err != nil {
		return false
	}
	u, err := url.Parse(token)
	if err != nil || u.Host == "" {
		return false
	}
	return true
}
type Result struct {
	Link   string
	Code   string
	Status string
}
func indexPage(w http.ResponseWriter, r *http.Request) {
	templ, _ := template.ParseFiles("index.html")
	result := Result{}
	if r.Method == "POST" {
		if !isValidUrl(r.FormValue("s")) {
			fmt.Println("���-�� �� ���")
			result.Status = "������������ ������!"
			result.Link = ""
		} else {
			result.Link = r.FormValue("s")
			result.Code = shorting()
			db, err := sql.Open("sqlite3", "myDB.db")
			if err != nil {
				panic(err)
			}
			defer db.Close()
			db.Exec("insert into links (link, short) values ($1, $2)", result.Link, result.Code)
			result.Status = "�������"
		}
	}
	templ.Execute(w, result)
}
func redirectTo(w http.ResponseWriter, r *http.Request) {
	var link string
	vars := mux.Vars(r)
	db, err := sql.Open("sqlite3", "myDB.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	rows := db.QueryRow("select link from links where short=$1 limit 1", vars["key"])
	rows.Scan(&link)
	fmt.Fprintf(w, "<script>location='%s';</script>", link)
}
func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", indexPage)
	router.HandleFunc("/s/{key}", redirectTo)
	log.Fatal(http.ListenAndServe(":8080", router))
}