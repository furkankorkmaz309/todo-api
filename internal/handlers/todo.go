package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/furkankorkmaz309/todo-api/internal/app"
	"github.com/furkankorkmaz309/todo-api/internal/models"
)

func GetTodos(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := `SELECT id, title, content, priority, created_at, due_date, done, category_id FROM todo`
		rows, err := app.DB.Query(query)
		if err != nil {
			respondError(w, app.ErrorLog, http.StatusInternalServerError, "Database error", err)
			return
		}
		defer rows.Close()

		var todos []models.Todo
		for rows.Next() {
			var todo models.Todo
			err = rows.Scan(&todo.ID, &todo.Title, &todo.Content, &todo.Priority, &todo.CreatedAt, &todo.DueDate, &todo.IsDone, &todo.CategoryID)
			if err != nil {
				respondError(w, app.ErrorLog, http.StatusInternalServerError, "Row could not read", err)
				return
			}
			todos = append(todos, todo)
		}

		respondJSON(w, http.StatusOK, todos, "Todos listed successfully.")
	}
}

func CreateTodo(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var todo models.Todo
		err := json.NewDecoder(r.Body).Decode(&todo)
		if err != nil {
			respondError(w, app.ErrorLog, http.StatusBadRequest, "Invalid JSON", err)
			return
		}

		if strings.TrimSpace(todo.Title) == "" {
			respondError(w, app.ErrorLog, http.StatusBadRequest, "Title is blank", fmt.Errorf("blank title"))
			return
		}
		if strings.TrimSpace(todo.Content) == "" {
			respondError(w, app.ErrorLog, http.StatusBadRequest, "Content is blank", fmt.Errorf("blank content"))
			return
		}
		if todo.Priority < 1 || todo.Priority > 5 {
			respondError(w, app.ErrorLog, http.StatusBadRequest, "Priority must be between 1-5", nil)
			return
		}

		todo.CreatedAt = time.Now()
		if todo.DueDate.Before(time.Now()) {
			respondError(w, app.ErrorLog, http.StatusBadRequest, "Due date can't be in the past", nil)
			return
		}
		todo.IsDone = false

		row := app.DB.QueryRow(`SELECT id FROM category WHERE id = ?`, todo.CategoryID)
		var tempID int
		err = row.Scan(&tempID)
		if err == sql.ErrNoRows {
			respondError(w, app.ErrorLog, http.StatusBadRequest, fmt.Sprintf("No category with ID %v", todo.CategoryID), nil)
			return
		}
		if err != nil {
			respondError(w, app.ErrorLog, http.StatusInternalServerError, "Database error", err)
			return
		}

		query := `INSERT INTO todo(title, content, priority, created_at, due_date, done, category_id) VALUES(?, ?, ?, ?, ?, ?, ?)`
		result, err := app.DB.Exec(query, todo.Title, todo.Content, todo.Priority, todo.CreatedAt, todo.DueDate, todo.IsDone, todo.CategoryID)
		if err != nil {
			respondError(w, app.ErrorLog, http.StatusInternalServerError, "Insert failed", err)
			return
		}
		id, err := result.LastInsertId()
		if err != nil {
			respondError(w, app.ErrorLog, http.StatusInternalServerError, "Failed to retrieve inserted ID", err)
			return
		}
		todo.ID = int(id)

		respondJSON(w, http.StatusCreated, todo, "Todo created successfully.")
	}
}

func GetTodo(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := takeIDFromURL(r, 2)
		if err != nil {
			respondError(w, app.ErrorLog, http.StatusBadRequest, "Invalid ID", err)
			return
		}

		var todo models.Todo
		query := `SELECT * FROM todo WHERE id = ?`
		row := app.DB.QueryRow(query, id)
		err = row.Scan(&todo.ID, &todo.Title, &todo.Content, &todo.Priority, &todo.CreatedAt, &todo.DueDate, &todo.IsDone, &todo.CategoryID)
		if err != nil {
			respondError(w, app.ErrorLog, http.StatusNotFound, fmt.Sprintf("No todo with ID %v", id), err)
			return
		}

		respondJSON(w, http.StatusOK, todo, "Todo fetched successfully.")
	}
}

func PatchTodo(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := takeIDFromURL(r, 2)
		if err != nil {
			respondError(w, app.ErrorLog, http.StatusBadRequest, "Invalid ID", err)
			return
		}

		var newTodo models.Todo
		err = json.NewDecoder(r.Body).Decode(&newTodo)
		if err != nil {
			respondError(w, app.ErrorLog, http.StatusBadRequest, "Invalid JSON body", err)
			return
		}

		var oldTodo models.Todo
		query := `SELECT id, title, content, priority, created_at, due_date, done, category_id FROM todo WHERE id = ?`
		row := app.DB.QueryRow(query, id)
		err = row.Scan(&oldTodo.ID, &oldTodo.Title, &oldTodo.Content, &oldTodo.Priority, &oldTodo.CreatedAt, &oldTodo.DueDate, &oldTodo.IsDone, &oldTodo.CategoryID)
		if err == sql.ErrNoRows {
			respondError(w, app.ErrorLog, http.StatusNotFound, fmt.Sprintf("No todo with ID %v", id), nil)
			return
		}
		if err != nil {
			respondError(w, app.ErrorLog, http.StatusInternalServerError, "Database error", err)
			return
		}

		responseString := "title, content, priority, due_date, is_done, category_id updated!"

		if strings.TrimSpace(newTodo.Title) != "" {
			oldTodo.Title = newTodo.Title
		} else {
			responseString = strings.ReplaceAll(responseString, "title, ", "")
		}
		if strings.TrimSpace(newTodo.Content) != "" {
			oldTodo.Content = newTodo.Content
		} else {
			responseString = strings.ReplaceAll(responseString, "content, ", "")
		}
		if newTodo.Priority >= 1 && newTodo.Priority <= 5 {
			oldTodo.Priority = newTodo.Priority
		} else {
			responseString = strings.ReplaceAll(responseString, "priority, ", "")
		}
		if newTodo.DueDate.After(time.Now()) {
			oldTodo.DueDate = newTodo.DueDate
		} else {
			responseString = strings.ReplaceAll(responseString, "due_date, ", "")
		}
		if oldTodo.IsDone == newTodo.IsDone {
			responseString = strings.ReplaceAll(responseString, "is_done, ", "")
		}
		oldTodo.IsDone = newTodo.IsDone
		if newTodo.CategoryID != 0 {
			oldTodo.CategoryID = newTodo.CategoryID
		} else {
			responseString = strings.ReplaceAll(responseString, "category_id ", "")
		}

		if responseString == "updated!" {
			respondError(w, app.ErrorLog, http.StatusBadRequest, "No fields provided for update", nil)
			return
		}

		queryUpdate := `UPDATE todo SET title = ?, content = ?, priority = ?, due_date = ?, done = ?, category_id = ? WHERE id = ?`
		result, err := app.DB.Exec(queryUpdate, oldTodo.Title, oldTodo.Content, oldTodo.Priority, oldTodo.DueDate, oldTodo.IsDone, oldTodo.CategoryID, id)
		if err != nil {
			respondError(w, app.ErrorLog, http.StatusInternalServerError, "Failed to update todo", err)
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			respondError(w, app.ErrorLog, http.StatusInternalServerError, "Could not get update result", err)
			return
		}
		if rowsAffected == 0 {
			respondError(w, app.ErrorLog, http.StatusNotFound, fmt.Sprintf("No todo with ID %v", id), nil)
			return
		}

		respondJSON(w, http.StatusOK, oldTodo, responseString)
	}
}

func DeleteTodo(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := takeIDFromURL(r, 2)
		if err != nil {
			respondError(w, app.ErrorLog, http.StatusBadRequest, "Invalid ID", err)
			return
		}

		query := `DELETE FROM todo WHERE id = ?`
		result, err := app.DB.Exec(query, id)
		if err != nil {
			respondError(w, app.ErrorLog, http.StatusInternalServerError, "Database error", err)
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			respondError(w, app.ErrorLog, http.StatusInternalServerError, "Could not retrieve delete result", err)
			return
		}

		if rowsAffected == 0 {
			respondError(w, app.ErrorLog, http.StatusNotFound, fmt.Sprintf("No todo with ID %v", id), nil)
			return
		}

		respondSuccess(w, http.StatusOK, fmt.Sprintf("Todo with ID %d deleted.", id))
	}
}
