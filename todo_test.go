package main

import "testing"

func TestCreateInsertsAtTop(t *testing.T) {
	setupTestDB(t)

	mustCreate(t, "first")
	mustCreate(t, "second")
	mustCreate(t, "third")

	todos, err := DB.GetTodos()
	if err != nil {
		t.Fatalf("GetTodos returned error: %s", err)
	}
	got := []string{}
	for _, todo := range todos {
		got = append(got, todo.Name)
	}
	want := []string{"third", "second", "first"}
	if len(got) != len(want) {
		t.Fatalf("expected %d todos, got %d (%v)", len(want), len(got), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("position %d: want %q, got %q (full order %v)", i, want[i], got[i], got)
		}
	}
}

func TestMarkCompleteHidesFromActiveList(t *testing.T) {
	setupTestDB(t)

	todo := mustCreate(t, "do the thing")
	if err := todo.MarkComplete(); err != nil {
		t.Fatalf("MarkComplete returned error: %s", err)
	}

	active, err := DB.GetTodos()
	if err != nil {
		t.Fatalf("GetTodos returned error: %s", err)
	}
	if len(active) != 0 {
		t.Errorf("expected completed todo to be hidden from active list, got %d items", len(active))
	}

	all, err := DB.GetAllTodos()
	if err != nil {
		t.Fatalf("GetAllTodos returned error: %s", err)
	}
	if len(all) != 1 {
		t.Fatalf("expected 1 todo in full list, got %d", len(all))
	}
	if !all[0].Completed {
		t.Errorf("expected todo to be marked completed")
	}
}

func TestDeleteRemovesFromAllLists(t *testing.T) {
	setupTestDB(t)

	todo := mustCreate(t, "delete me")
	if err := todo.Delete(); err != nil {
		t.Fatalf("Delete returned error: %s", err)
	}

	all, err := DB.GetAllTodos()
	if err != nil {
		t.Fatalf("GetAllTodos returned error: %s", err)
	}
	if len(all) != 0 {
		t.Errorf("expected deleted todo to be excluded, got %d items", len(all))
	}
}

func TestUpdateContent(t *testing.T) {
	setupTestDB(t)

	todo := mustCreate(t, "with content")
	content := "some detailed notes"
	todo.Content = &content
	if err := todo.UpdateContent(); err != nil {
		t.Fatalf("UpdateContent returned error: %s", err)
	}

	got, err := DB.GetTodo(todo.Id)
	if err != nil {
		t.Fatalf("GetTodo returned error: %s", err)
	}
	if got.Content == nil {
		t.Fatalf("expected content to be set, got nil")
	}
	if *got.Content != content {
		t.Errorf("want content %q, got %q", content, *got.Content)
	}
}

func TestGetTodoNotFound(t *testing.T) {
	setupTestDB(t)

	if _, err := DB.GetTodo(12345); err == nil {
		t.Errorf("expected error for missing todo, got nil")
	}
}
