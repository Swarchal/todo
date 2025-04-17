package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	con *sql.DB
}

// create datbase connection, create database table if needed
func CreateDb() Database {
	const createStr string = `
    CREATE TABLE IF NOT EXISTS Todo (
      id        INTEGER  PRIMARY KEY,
      name      STRING   NOT NULL,
      completed INTEGER  NOT NULL DEFAULT 0,
      date      DATETIME DEFAULT CURRENT_TIMESTAMP,
      deleted   INTEGER  NOT NULL DEFAULT 0,
      ordering  INTEGER  NOT NULL DEFAULT 0
    )
    `
	db, err := sql.Open("sqlite3", "./db.sqlite")
	if err != nil {
		log.Fatalln(err)
	}
	_, err = db.Exec(createStr)
	if err != nil {
		log.Fatalln(err)
	}
	return Database{db}
}

// generic query Todo items
func (db Database) queryRows(query string) ([]Todo, error) {
	rows, err := db.con.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var todoList []Todo
	for rows.Next() {
		var t Todo
		err := rows.Scan(
			&t.Id, &t.Name, &t.Completed, &t.Date, &t.Deleted, &t.Ordering,
		)
		if err != nil {
			return todoList, err
		}
		todoList = append(todoList, t)
	}
	if err = rows.Err(); err != nil {
		return todoList, err
	}
	return todoList, nil
}

// get all todo items (including those done) but excluding deleted
func (db Database) GetAllTodos() ([]Todo, error) {
	query := "SELECT id, name, completed, date, deleted, ordering FROM Todo WHERE deleted != 1 ORDER BY ordering"
	return db.queryRows(query)
}

// get todo items which are not done
func (db Database) GetTodos() ([]Todo, error) {
	query := "SELECT id, name, completed, date, deleted, ordering FROM Todo WHERE deleted != 1 AND completed != 1 ORDER BY ordering"
	return db.queryRows(query)
}

// get todo from database by id
func (db Database) GetTodo(id int64) (Todo, error) {
	row := db.con.QueryRow("SELECT id, name, completed, date, deleted, ordering FROM Todo WHERE id = ?", id)
	var t Todo
	err := row.Scan(&t.Id, &t.Name, &t.Completed, &t.Date, &t.Deleted, &t.Ordering)
	if err != nil {
		return t, err
	}
	return t, nil
}

func (db Database) CompleteTodo(id int64) error {
	_, err := db.con.Exec("UPDATE Todo SET completed = 1 WHERE id = ?", id)
	return err
}
