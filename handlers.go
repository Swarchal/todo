package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func handleIndex(w http.ResponseWriter, r *http.Request) {
	todos, _ := DB.GetTodos()
	t, _ := template.ParseFiles("templates/index.html")

	t.Execute(w, todos)
}

func handleAdd(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		err := r.ParseForm()
		if err != nil {
			log.Fatal(err)
		}
		name := r.Form.Get("Name")
		t := Todo{Name: name}
		err = t.Save()
		if err != nil {
			log.Fatalf("Cannot save TODO item: %s", err)
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func handleComplete(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/complete/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Fatalf("Failed to convert %s to an int\n", idStr)
	}
	todo, err := DB.GetTodo(id)
	if err != nil {
		log.Fatalf("Cannot find TODO id:%b\n", id)
	}
	err = todo.MarkComplete()
	if err != nil {
		log.Fatalf("Cannot complete TODO id %b: %s\n", id, err)
	}
}

func handleSortItems(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Fatalf("Failed to parse form: %s\n", err)
	}
	idStrs := r.Form["item"]
	// convert slice from strings to int64s
	ids := make([]int64, len(idStrs))
	for idx, idStr := range idStrs {
		id32, _ := strconv.Atoi(idStr)
		ids[idx] = int64(id32)
	}
	// set ordering value for each TODO.Id
	// save in database
	for idx, id := range ids {
		t, _ := DB.GetTodo(id)
		t.Ordering = idx + 1
		err := t.Update()
		if err != nil {
			log.Printf("Failed to update TODO id: %b: %s", t.Id, err)
		}
	}
	fmt.Fprintf(w, "")
}
