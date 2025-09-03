package main

import (
	"fmt"
	"log"
	"net/http"
)

var DB Database

func main() {
	DB = CreateDb()
	// Have to serve static css files in Go, otherwise css file is not found.
	// https://www.devasking.com/issue/how-to-run-html-with-css-using-golang
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))

	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/add", handleAdd)
	http.HandleFunc("/complete/", handleComplete)
	http.HandleFunc("/delete/", handleDelete)
	http.HandleFunc("/sort-items", handleSortItems)
	http.HandleFunc("/detail/", handleGetDetail)
	http.HandleFunc("/update-content/", handleUpdateContent)

	fmt.Println("Running http server...")
	if err := http.ListenAndServe("127.0.0.1:3333", nil); err != nil {
		log.Fatal(err)
	}
}
