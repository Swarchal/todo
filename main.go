package main

import (
	"fmt"
	"log"
	"net/http"
)

var DB Database

func main() {
	DB = CreateDb()

	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/add", handleAdd)
	http.HandleFunc("/complete/", handleComplete)
	http.HandleFunc("/delete/", handleDelete)
	http.HandleFunc("/sort-items", handleSortItems)

	fmt.Println("Running http server...")
	if err := http.ListenAndServe("127.0.0.1:3333", nil); err != nil {
		log.Fatal(err)
	}
}
