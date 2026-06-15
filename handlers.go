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
	todos, err := DB.GetTodos()
	if err != nil {
		log.Printf("Cannot fetch TODOs: %s\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Printf("Cannot parse template: %s\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if err := t.Execute(w, todos); err != nil {
		log.Printf("Cannot render template: %s\n", err)
	}
}

func handleAdd(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		err := r.ParseForm()
		if err != nil {
			log.Printf("Failed to parse form: %s\n", err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		name := r.Form.Get("Name")
		t := Todo{Name: name}
		err = t.Create()
		if err != nil {
			log.Printf("Cannot save TODO item: %s\n", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func handleComplete(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		idStr := strings.TrimPrefix(r.URL.Path, "/complete/")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Printf("Failed to convert %s to an int\n", idStr)
			http.Error(w, "Invalid TODO id", http.StatusBadRequest)
			return
		}
		todo, err := DB.GetTodo(id)
		if err != nil {
			log.Printf("Cannot find TODO id:%d\n", id)
			http.Error(w, "TODO not found", http.StatusNotFound)
			return
		}
		err = todo.MarkComplete()
		if err != nil {
			log.Printf("Cannot complete TODO id %d: %s\n", id, err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		idStr := strings.TrimPrefix(r.URL.Path, "/delete/")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Printf("Failed to convert %s to an int\n", idStr)
			http.Error(w, "Invalid TODO id", http.StatusBadRequest)
			return
		}
		todo, err := DB.GetTodo(id)
		if err != nil {
			log.Printf("Cannot find TODO id:%d\n", id)
			http.Error(w, "TODO not found", http.StatusNotFound)
			return
		}
		err = todo.Delete()
		if err != nil {
			log.Printf("Cannot delete TODO id %d: %s\n", id, err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}

func handleSortItems(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			log.Printf("Failed to parse form: %s\n", err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
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
				log.Printf("Failed to update TODO id: %d: %s\n", t.Id, err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		}
	}
	fmt.Fprintf(w, "")
}

func handleGetDetail(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		idStr := strings.TrimPrefix(r.URL.Path, "/detail/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Printf("Failed to convert %s to an int\n", idStr)
			http.Error(w, "Invalid TODO id", http.StatusBadRequest)
			return
		}
		id64 := int64(id)
		todo, err := DB.GetTodo(id64)
		if err != nil {
			log.Printf("Cannot find TODO id:%d err: %s\n", id64, err.Error())
			http.Error(w, "TODO not found", http.StatusNotFound)
			return
		}
		t, err := template.ParseFiles("templates/detail.html")
		if err != nil {
			log.Printf("Cannot parse template: %s\n", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		t.Execute(w, todo)
	}
}

func handleUpdateContent(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			log.Printf("Failed to parse form: %s\n", err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		idStr := r.Form.Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Printf("Failed to convert %s to an int\n", idStr)
			http.Error(w, "Invalid TODO id", http.StatusBadRequest)
			return
		}
		id64 := int64(id)
		todo, err := DB.GetTodo(id64)
		if err != nil {
			log.Printf("Cannot find TODO id:%d err: %s\n", id64, err.Error())
			http.Error(w, "TODO not found", http.StatusNotFound)
			return
		}
		content := r.Form.Get("content")
		todo.Content = &content
		err = todo.UpdateContent()
		if err != nil {
			log.Printf("Cannot update content for TODO id:%d err: %s\n", id64, err.Error())
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
