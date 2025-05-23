package models

import "time"

type Todo struct {
	ID         int       `json:"id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	Priority   int       `json:"priority"`
	CreatedAt  time.Time `json:"created_at"`
	DueDate    time.Time `json:"due_date"`
	IsDone     bool      `json:"is_done"`
	CategoryID int       `json:"category_id"`
}
