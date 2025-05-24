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
			fail(w, app.ErrorLog, "GetTodos", err, "Database error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var todos []models.Todo

		for rows.Next() {
			var todo models.Todo
			err = rows.Scan(&todo.ID, &todo.Title, &todo.Content, &todo.Priority, &todo.CreatedAt, &todo.DueDate, &todo.IsDone, &todo.CategoryID)
			if err != nil {
				fail(w, app.ErrorLog, "GetTodos", err, "Row could not read", http.StatusInternalServerError)
				return
			}
			todos = append(todos, todo)
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(todos)
		if err != nil {
			fail(w, app.ErrorLog, "GetTodos", err, "JSON encoding error", http.StatusInternalServerError)
			return
		}
	}
}

func CreateTodo(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var todo models.Todo
		err := json.NewDecoder(r.Body).Decode(&todo)
		if err != nil {
			fail(w, app.ErrorLog, "CreateTodo", err, err.Error(), http.StatusBadRequest)
			return
		}

		if strings.TrimSpace(todo.Title) == "" {
			err := fmt.Errorf("title field is blank")
			fail(w, app.ErrorLog, "CreateTodo", err, err.Error(), http.StatusBadRequest)
			return
		}
		if strings.TrimSpace(todo.Content) == "" {
			err := fmt.Errorf("content field is blank")
			fail(w, app.ErrorLog, "CreateTodo", err, err.Error(), http.StatusBadRequest)
			return
		}
		if todo.Priority < 1 || todo.Priority > 5 {
			err := fmt.Errorf("priority should be in range 1 - 5")
			fail(w, app.ErrorLog, "CreateTodo", err, err.Error(), http.StatusBadRequest)
			return
		}

		todo.CreatedAt = time.Now()
		if todo.DueDate.Before(time.Now()) {
			err := fmt.Errorf("due date can not be in past")
			fail(w, app.ErrorLog, "CreateTodo", err, err.Error(), http.StatusBadRequest)
			return
		}

		todo.IsDone = false

		isCategoryExistQuery := `SELECT id FROM category WHERE id = ?`
		row := app.DB.QueryRow(isCategoryExistQuery, todo.CategoryID)
		var temporaryID int
		err = row.Scan(&temporaryID)
		if err == sql.ErrNoRows {
			err := fmt.Errorf("there is no category with ID %v", todo.CategoryID)
			fail(w, app.ErrorLog, "CreateTodo", err, err.Error(), http.StatusBadRequest)
			return
		}
		if err != nil {
			fail(w, app.ErrorLog, "CreateTodo", err, "Database error", http.StatusInternalServerError)
			return
		}

		query := `INSERT INTO todo(title, content, priority, created_at, due_date, done, category_id) VALUES(?, ?, ?, ?, ?, ?, ?)`
		result, err := app.DB.Exec(query, todo.Title, todo.Content, todo.Priority, todo.CreatedAt, todo.DueDate, todo.IsDone, todo.CategoryID)
		if err != nil {
			fail(w, app.ErrorLog, "CreateTodo", err, "Database error", http.StatusInternalServerError)
			return
		}
		id, err := result.LastInsertId()
		if err != nil {
			fail(w, app.ErrorLog, "CreateTodo", err, "Failed to retrieve inserted ID", http.StatusInternalServerError)
			return
		}
		todo.ID = int(id)

		// burada çakışma olabilir de bunu çözecek kafa yok bende
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(todo)
		if err != nil {
			fail(w, app.ErrorLog, "CreateTodo", err, "JSON encode error", http.StatusInternalServerError)
			return
		}
	}
}

func GetTodo(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := takeIDFromURL(r, 2)
		if err != nil {
			fail(w, app.ErrorLog, "GetTodo", err, err.Error(), http.StatusBadRequest)
			return
		}

		var todo models.Todo
		query := `SELECT * FROM todo WHERE id = ?`
		row := app.DB.QueryRow(query, id)
		err = row.Scan(&todo.ID, &todo.Title, &todo.Content, &todo.Priority, &todo.CreatedAt, &todo.DueDate, &todo.IsDone, &todo.CategoryID)
		if err != nil {
			err := fmt.Errorf("there is no todo with ID %v", id)
			fail(w, app.ErrorLog, "GetTodo", err, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(todo)
	}
}

func PatchTodo(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := takeIDFromURL(r, 2)
		if err != nil {
			fail(w, app.ErrorLog, "UpdateTodo", err, err.Error(), http.StatusBadRequest)
			return
		}

		var newTodo models.Todo
		err = json.NewDecoder(r.Body).Decode(&newTodo)
		if err != nil {
			fail(w, app.ErrorLog, "UpdateTodo", err, "Invalid JSON body", http.StatusBadRequest)
			return
		}

		var oldTodo models.Todo
		query := `SELECT id, title, content, priority, created_at, due_date, done, category_id FROM todo WHERE id = ?`
		row := app.DB.QueryRow(query, id)
		err = row.Scan(&oldTodo.ID, &oldTodo.Title, &oldTodo.Content, &oldTodo.Priority, &oldTodo.CreatedAt, &oldTodo.DueDate, &oldTodo.IsDone, &oldTodo.CategoryID)
		if err == sql.ErrNoRows {
			err := fmt.Errorf("there is no todo with ID %v", id)
			fail(w, app.ErrorLog, "UpdateTodo", err, err.Error(), http.StatusNotFound)
			return
		}
		if err != nil {
			err := fmt.Errorf("database error: %v", err)
			fail(w, app.ErrorLog, "UpdateTodo", err, err.Error(), http.StatusInternalServerError)
			return
		}

		responseString := "title, content, priority, due_date, is_done, category_id updated!"

		if strings.TrimSpace(newTodo.Title) == "" {
			newTodo.Title = oldTodo.Title
		} else {
			responseString = strings.ReplaceAll(responseString, "title, ", "")
		}
		if strings.TrimSpace(newTodo.Content) == "" {
			newTodo.Content = oldTodo.Content
		} else {
			responseString = strings.ReplaceAll(responseString, "content, ", "")
		}
		if newTodo.Priority < 1 || newTodo.Priority > 5 {
			err := fmt.Errorf("priority should be in range 1 - 5")
			fail(w, app.ErrorLog, "UpdateTodo", err, err.Error(), http.StatusBadRequest)
			return
		} else {
			responseString = strings.ReplaceAll(responseString, "priority, ", "")
		}
		if newTodo.Priority == 0 {
			newTodo.Priority = oldTodo.Priority
		}
		if newTodo.DueDate.IsZero() {
			newTodo.DueDate = oldTodo.DueDate
		} else if newTodo.DueDate.Before(time.Now()) {
			err := fmt.Errorf("due date can not be in past")
			fail(w, app.ErrorLog, "UpdateTodo", err, err.Error(), http.StatusBadRequest)
			return
		} else {
			responseString = strings.ReplaceAll(responseString, "due_date, ", "")
		}
		if oldTodo.IsDone == newTodo.IsDone {
			responseString = strings.ReplaceAll(responseString, "is_done, ", "")
		}
		oldTodo.IsDone = newTodo.IsDone
		if newTodo.CategoryID == 0 {
			newTodo.CategoryID = oldTodo.CategoryID
		} else {
			responseString = strings.ReplaceAll(responseString, "category_id", "")
		}

		queryUpdate := `UPDATE todo SET title = ?, content = ?, priority = ?, due_date = ?, done = ?, category_id = ? WHERE id = ?`
		result, err := app.DB.Exec(queryUpdate, newTodo.Title, newTodo.Content, newTodo.Priority, newTodo.DueDate, newTodo.IsDone, newTodo.CategoryID, id)
		if err != nil {
			fail(w, app.ErrorLog, "UpdateTodo", err, "Failed to update todo", http.StatusInternalServerError)
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			fail(w, app.ErrorLog, "UpdateTodo", err, "Result could not take", http.StatusInternalServerError)
			return
		}

		if rowsAffected == 0 {
			err := fmt.Errorf("there is no todo with ID %v", id)
			fail(w, app.ErrorLog, "UpdateTodo", err, err.Error(), http.StatusNotFound)
			return
		}

		info := fmt.Sprintf("Todo with ID %v updated!", id)
		infoMap := map[string]any{
			"message":  info,
			"response": responseString,
			"data":     newTodo,
		}
		json.NewEncoder(w).Encode(infoMap)
	}
}

func DeleteTodo(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := takeIDFromURL(r, 2)
		if err != nil {
			fail(w, app.ErrorLog, "DeleteTodo", err, err.Error(), http.StatusBadRequest)
			return
		}

		query := `DELETE FROM todo WHERE id = ?`
		result, err := app.DB.Exec(query, id)
		if err != nil {
			fail(w, app.ErrorLog, "DeleteTodo", err, "Database error", http.StatusInternalServerError)
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			fail(w, app.ErrorLog, "DeleteTodo", err, "Result could not take", http.StatusInternalServerError)
			return
		}

		if rowsAffected == 0 {
			err := fmt.Errorf("there is no todo with ID %v", id)
			fail(w, app.ErrorLog, "DeleteTodo", err, err.Error(), http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
