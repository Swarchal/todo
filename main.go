package main

import (
	"fmt"
	"net/http"
)

var DB Database

func main() {
	DB = CreateDb()

	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/add", handleAdd)
	http.HandleFunc("/complete/", handleComplete)

	fmt.Println("Running http server...")
	http.ListenAndServe(":3333", nil)
}
