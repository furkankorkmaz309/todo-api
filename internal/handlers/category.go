package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/furkankorkmaz309/todo-api/internal/app"
	"github.com/furkankorkmaz309/todo-api/internal/models"
)

func GetCategories(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		query := `SELECT id, name, description FROM category`
		rows, err := app.DB.Query(query)
		if err != nil {
			fail(w, app.ErrorLog, "GetCategories", err, "Database error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var categories []models.Category

		for rows.Next() {
			var category models.Category
			err := rows.Scan(&category.ID, &category.Name, &category.Description)
			if err != nil {
				fail(w, app.ErrorLog, "GetCategories", err, "Row could not read", http.StatusInternalServerError)
				return
			}
			categories = append(categories, category)
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(categories)
		if err != nil {
			fail(w, app.ErrorLog, "GetCategories", err, "JSON encoding error", http.StatusInternalServerError)
			return
		}
	}
}

func AddCategory(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var input models.Category
		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			fail(w, app.ErrorLog, "AddCategory", err, "Invalid JSOn body", http.StatusBadRequest)
			return
		}

		// buraya errorlog yaz
		if strings.TrimSpace(input.Name) == "" {
			err := fmt.Errorf("name field is blank")
			fail(w, app.ErrorLog, "AddCategory", err, err.Error(), http.StatusBadRequest)
			return
		}

		if len(input.Name) > 30 {
			err := fmt.Errorf("name field is too long")
			fail(w, app.ErrorLog, "AddCategory", err, err.Error(), http.StatusBadRequest)
			return
		}

		if len(input.Description) > 100 {
			err := fmt.Errorf("description field is too long")
			fail(w, app.ErrorLog, "AddCategory", err, err.Error(), http.StatusBadRequest)
			return
		}

		query := `INSERT INTO category (name, description) VALUES (?,?)`
		result, err := app.DB.Exec(query, input.Name, input.Description)
		if err != nil {
			fail(w, app.ErrorLog, "AddCategory", err, "Database error", http.StatusInternalServerError)
			return
		}

		id, err := result.LastInsertId()
		if err != nil {
			fail(w, app.ErrorLog, "AddCategory", err, "Failed to retrieve inserted ID", http.StatusInternalServerError)
			return
		}
		input.ID = int(id)

		// burada çakışma olabilir de bunu çözecek kafa yok bende
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(input)
		if err != nil {
			fail(w, app.ErrorLog, "AddCategory", err, "JSON encoding error", http.StatusInternalServerError)
			return
		}
	}
}

func UpdateCategory(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id, err := takeIDFromURL(r, 2)
		if err != nil {
			fail(w, app.ErrorLog, "UpdateCategory", err, err.Error(), http.StatusBadRequest)
			return
		}

		query := `SELECT * FROM category WHERE id = ?`
		row := app.DB.QueryRow(query, id)

		var oldCategory models.Category
		err = row.Scan(&oldCategory.ID, &oldCategory.Name, &oldCategory.Description)
		if err != nil {
			fail(w, app.ErrorLog, "UpdateCategory", err, "Row could not read", http.StatusInternalServerError)
			return
		}

		var newCategory models.Category
		err = json.NewDecoder(r.Body).Decode(&newCategory)
		if err != nil {
			fail(w, app.ErrorLog, "UpdateCategory", err, "Invalid JSON body", http.StatusBadRequest)
			return
		}

		if len(newCategory.Name) > 30 {
			err := fmt.Errorf("name field is too long")
			fail(w, app.ErrorLog, "UpdateCategory", err, err.Error(), http.StatusBadRequest)
			return
		}
		if len(newCategory.Description) > 100 {
			err := fmt.Errorf("description field is too long")
			fail(w, app.ErrorLog, "UpdateCategory", err, err.Error(), http.StatusBadRequest)
			return
		}
		if strings.TrimSpace(newCategory.Name) == "" {
			newCategory.Name = oldCategory.Name
		}
		if strings.TrimSpace(newCategory.Description) == "" {
			newCategory.Description = oldCategory.Description
		}

		queryUpdate := `UPDATE category SET name = ?, description = ? WHERE id = ?`
		result, err := app.DB.Exec(queryUpdate, newCategory.Name, newCategory.Description, id)
		if err != nil {
			fail(w, app.ErrorLog, "UpdateCategory", err, "Failed to update category", http.StatusInternalServerError)
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			fail(w, app.ErrorLog, "UpdateCategory", err, "Result not take", http.StatusInternalServerError)
			return
		}

		if rowsAffected == 0 {
			err := fmt.Errorf("there is no category with ID %v", id)
			fail(w, app.ErrorLog, "UpdateCategory", err, err.Error(), http.StatusNotFound)
			return
		}

		info := fmt.Sprintf("Category with ID %v updated!", id)
		infoMap := map[string]string{
			"message": info,
		}
		json.NewEncoder(w).Encode(infoMap)
	}
}

func DeleteCategory(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id, err := takeIDFromURL(r, 2)
		if err != nil {
			fail(w, app.ErrorLog, "DeleteCategory", err, err.Error(), http.StatusBadRequest)
			return
		}

		query := `DELETE FROM category WHERE id = ?`
		result, err := app.DB.Exec(query, id)
		if err != nil {
			fail(w, app.ErrorLog, "DeleteCategory", err, "Database error", http.StatusInternalServerError)
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			fail(w, app.ErrorLog, "DeleteCategory", err, "Result not take", http.StatusInternalServerError)
			return
		}

		if rowsAffected == 0 {
			err := fmt.Errorf("there is no category with ID %v", id)
			fail(w, app.ErrorLog, "DeleteCategory", err, err.Error(), http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
