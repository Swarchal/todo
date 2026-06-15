package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
)

// postForm builds a POST request whose body is the url-encoded form values.
func postForm(path string, form url.Values) *http.Request {
	req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req
}

func TestHandleAddCreatesTodo(t *testing.T) {
	setupTestDB(t)

	w := httptest.NewRecorder()
	handleAdd(w, postForm("/add", url.Values{"Name": {"buy milk"}}))

	if w.Code != http.StatusSeeOther {
		t.Errorf("want status %d, got %d", http.StatusSeeOther, w.Code)
	}
	todos, err := DB.GetTodos()
	if err != nil {
		t.Fatalf("GetTodos returned error: %s", err)
	}
	if len(todos) != 1 || todos[0].Name != "buy milk" {
		t.Errorf("expected one todo named %q, got %+v", "buy milk", todos)
	}
}

func TestHandleAddRejectsGet(t *testing.T) {
	setupTestDB(t)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/add", nil)
	handleAdd(w, req)

	if w.Code != http.StatusSeeOther {
		t.Errorf("want redirect %d for GET, got %d", http.StatusSeeOther, w.Code)
	}
	todos, _ := DB.GetTodos()
	if len(todos) != 0 {
		t.Errorf("GET /add should not create a todo, got %d", len(todos))
	}
}

func TestHandleComplete(t *testing.T) {
	setupTestDB(t)
	todo := mustCreate(t, "finish me")

	w := httptest.NewRecorder()
	handleComplete(w, httptest.NewRequest(http.MethodPost, "/complete/"+strconv.FormatInt(todo.Id, 10), nil))

	got, err := DB.GetTodo(todo.Id)
	if err != nil {
		t.Fatalf("GetTodo returned error: %s", err)
	}
	if !got.Completed {
		t.Errorf("expected todo to be completed")
	}
}

func TestHandleCompleteInvalidIdDoesNotCrash(t *testing.T) {
	setupTestDB(t)

	w := httptest.NewRecorder()
	handleComplete(w, httptest.NewRequest(http.MethodPost, "/complete/not-a-number", nil))

	if w.Code != http.StatusBadRequest {
		t.Errorf("want status %d for bad id, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestHandleCompleteMissingIdReturns404(t *testing.T) {
	setupTestDB(t)

	w := httptest.NewRecorder()
	handleComplete(w, httptest.NewRequest(http.MethodPost, "/complete/9999", nil))

	if w.Code != http.StatusNotFound {
		t.Errorf("want status %d for missing id, got %d", http.StatusNotFound, w.Code)
	}
}

func TestHandleDelete(t *testing.T) {
	setupTestDB(t)
	todo := mustCreate(t, "remove me")

	w := httptest.NewRecorder()
	handleDelete(w, httptest.NewRequest(http.MethodPost, "/delete/"+strconv.FormatInt(todo.Id, 10), nil))

	all, err := DB.GetAllTodos()
	if err != nil {
		t.Fatalf("GetAllTodos returned error: %s", err)
	}
	if len(all) != 0 {
		t.Errorf("expected todo to be deleted, got %d items", len(all))
	}
}

func TestHandleIndexRendersTodos(t *testing.T) {
	setupTestDB(t)
	mustCreate(t, "visible todo")

	w := httptest.NewRecorder()
	handleIndex(w, httptest.NewRequest(http.MethodGet, "/", nil))

	if w.Code != http.StatusOK {
		t.Errorf("want status %d, got %d", http.StatusOK, w.Code)
	}
	if !strings.Contains(w.Body.String(), "visible todo") {
		t.Errorf("expected rendered page to contain todo name, body: %s", w.Body.String())
	}
}

func TestHandleSortItems(t *testing.T) {
	setupTestDB(t)
	a := mustCreate(t, "a")
	b := mustCreate(t, "b")
	c := mustCreate(t, "c")

	// Submit a new explicit order: b, c, a
	form := url.Values{"item": {
		strconv.FormatInt(b.Id, 10),
		strconv.FormatInt(c.Id, 10),
		strconv.FormatInt(a.Id, 10),
	}}
	w := httptest.NewRecorder()
	handleSortItems(w, postForm("/sort-items", form))

	todos, err := DB.GetTodos()
	if err != nil {
		t.Fatalf("GetTodos returned error: %s", err)
	}
	got := []string{}
	for _, todo := range todos {
		got = append(got, todo.Name)
	}
	want := []string{"b", "c", "a"}
	for i := range want {
		if i >= len(got) || got[i] != want[i] {
			t.Fatalf("want order %v, got %v", want, got)
		}
	}
}

func TestHandleUpdateContent(t *testing.T) {
	setupTestDB(t)
	todo := mustCreate(t, "note holder")

	form := url.Values{
		"id":      {strconv.FormatInt(todo.Id, 10)},
		"content": {"remember the eggs"},
	}
	w := httptest.NewRecorder()
	handleUpdateContent(w, postForm("/update-content/", form))

	got, err := DB.GetTodo(todo.Id)
	if err != nil {
		t.Fatalf("GetTodo returned error: %s", err)
	}
	if got.Content == nil || *got.Content != "remember the eggs" {
		t.Errorf("expected content to be saved, got %v", got.Content)
	}
}
