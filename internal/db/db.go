package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "../../internal/data/todo-api.db")
	if err != nil {
		return nil, fmt.Errorf("an error occured while opening database : %v", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("an error occured while connecting database : %v", err)
	}

	queryCategory := `CREATE TABLE IF NOT EXISTS category (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT,
	description TEXT
	)`

	_, err = db.Exec(queryCategory)
	if err != nil {
		return nil, fmt.Errorf("an error occurred while creating category table : %v", err)
	}

	queryTodo := `CREATE TABLE IF NOT EXISTS todo (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	title TEXT,
	content TEXT,
	priority INTEGER,
	created_at TIMESTAMP,
	due_date TIMESTAMP,
	done BOOLEAN DEFAULT 0,
	category_id INT,
	FOREIGN KEY (category_id) REFERENCES category(id)
	)`

	_, err = db.Exec(queryTodo)
	if err != nil {
		return nil, fmt.Errorf("an error occurred while creating todo table : %v", err)
	}
	return db, nil
}
