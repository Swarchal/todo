package main

import (
	"time"
)

type Todo struct {
	Id        int64
	Name      string
	Completed bool
	Date      time.Time
	Deleted   bool
	Ordering  int
}

func (t *Todo) Save() error {
	res, err := DB.con.Exec("INSERT INTO Todo (name) VALUES (?)", t.Name)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	t.Id = id
	return err
}

func (t *Todo) Delete() error {
	_, err := DB.con.Exec("UPDATE Todo SET deleted = 1 WHERE id = ?", t.Id)
	t.Deleted = true
	return err
}

func (t *Todo) MarkComplete() error {
	_, err := DB.con.Exec("UPDATE Todo SET completed = 1 WHERE id = ?", t.Id)
	t.Completed = true
	return err
}

func (t *Todo) Update() error {
	_, err := DB.con.Exec(
		"UPDATE Todo SET name = ?, completed = ?, deleted = ?, ordering = ? WHERE id = ?",
		t.Name, t.Completed, t.Deleted, t.Ordering, t.Id,
	)
	return err
}
