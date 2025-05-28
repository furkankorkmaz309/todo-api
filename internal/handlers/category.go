package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/furkankorkmaz309/todo-api/internal/app"
	"github.com/furkankorkmaz309/todo-api/internal/models"
	"github.com/go-chi/chi"
)

func GetCategories(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		query := `SELECT id, name, description FROM category`
		rows, err := app.DB.Query(query)
		if err != nil {
			respondError(w, app.ErrorLog, http.StatusInternalServerError, "Database error", err)
			return
		}
		defer rows.Close()

		var categories []models.Category

		for rows.Next() {
			var category models.Category
			err := rows.Scan(&category.ID, &category.Name, &category.Description)
			if err != nil {
				respondError(w, app.ErrorLog, http.StatusInternalServerError, "Row scan error", err)
				return
			}
			categories = append(categories, category)
		}

		respondJSON(w, http.StatusOK, categories, "Categories listed successfully!")
	}
}

func AddCategory(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var input models.Category
		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			respondError(w, app.ErrorLog, http.StatusBadRequest, "Invalid JSON body", err)
			return
		}

		if strings.TrimSpace(input.Name) == "" {
			respondError(w, app.ErrorLog, http.StatusBadRequest, "Name field is blank", fmt.Errorf("blank name"))
			return
		}

		if len(input.Name) > 30 {
			respondError(w, app.ErrorLog, http.StatusBadRequest, "Name field is too long", nil)
			return
		}

		if len(input.Description) > 100 {
			respondError(w, app.ErrorLog, http.StatusBadRequest, "Description field is too long", nil)
			return
		}

		query := `INSERT INTO category (name, description) VALUES (?,?)`
		result, err := app.DB.Exec(query, input.Name, input.Description)
		if err != nil {
			respondError(w, app.ErrorLog, http.StatusInternalServerError, "Database error", err)
			return
		}

		id, err := result.LastInsertId()
		if err != nil {
			respondError(w, app.ErrorLog, http.StatusInternalServerError, "Failed to retrieve inserted ID", err)
			return
		}
		input.ID = int(id)

		respondJSON(w, http.StatusCreated, input, "Category created successfully!")
	}
}

func PatchCategory(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			err = fmt.Errorf("an error occurred while converting string to integer: %v", err)
			respondError(w, app.ErrorLog, http.StatusBadRequest, "Invalid category ID", err)
			return
		}

		query := `SELECT * FROM category WHERE id = ?`
		row := app.DB.QueryRow(query, id)

		var oldCategory models.Category
		err = row.Scan(&oldCategory.ID, &oldCategory.Name, &oldCategory.Description)
		if err != nil {
			respondError(w, app.ErrorLog, http.StatusNotFound, "Category not found", err)
			return
		}

		var newCategory models.Category
		err = json.NewDecoder(r.Body).Decode(&newCategory)
		if err != nil {
			respondError(w, app.ErrorLog, http.StatusBadRequest, "Invalid JSON body", err)
			return
		}

		responseString := "name, description updated!"

		if len(newCategory.Name) > 30 {
			respondError(w, app.ErrorLog, http.StatusBadRequest, "Name is too long", nil)
			return
		}
		if len(newCategory.Description) > 100 {
			respondError(w, app.ErrorLog, http.StatusBadRequest, "Description is too long", nil)
			return
		}
		if strings.TrimSpace(newCategory.Name) == "" {
			newCategory.Name = oldCategory.Name
			responseString = strings.ReplaceAll(responseString, "name, ", "")
		}
		if strings.TrimSpace(newCategory.Description) == "" {
			newCategory.Description = oldCategory.Description
			responseString = strings.ReplaceAll(responseString, "description ", "")
		}

		if responseString == "updated!" {
			respondError(w, app.ErrorLog, http.StatusBadRequest, "No fields provided for update", nil)
			return
		}

		queryUpdate := `UPDATE category SET name = ?, description = ? WHERE id = ?`
		result, err := app.DB.Exec(queryUpdate, newCategory.Name, newCategory.Description, id)
		if err != nil {
			respondError(w, app.ErrorLog, http.StatusInternalServerError, "Failed to update category", err)
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			respondError(w, app.ErrorLog, http.StatusInternalServerError, "Could not retrieve update result", err)
			return
		}

		if rowsAffected == 0 {
			respondError(w, app.ErrorLog, http.StatusNotFound, fmt.Sprintf("No category with ID %d", id), nil)
			return
		}

		newCategory.ID = id

		respondJSON(w, http.StatusOK, newCategory, responseString)
	}
}

func DeleteCategory(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			err = fmt.Errorf("an error occurred while converting string to integer: %v", err)
			respondError(w, app.ErrorLog, http.StatusBadRequest, "Invalid category ID", err)
			return
		}

		query := `DELETE FROM category WHERE id = ?`
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
			respondError(w, app.ErrorLog, http.StatusNotFound, fmt.Sprintf("No category with ID %d", id), nil)
			return
		}

		respondSuccess(w, http.StatusOK, fmt.Sprintf("Category with ID %d deleted.", id))
	}
}
