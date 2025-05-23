package routes

import (
	"net/http"

	"github.com/furkankorkmaz309/todo-api/internal/app"
	"github.com/furkankorkmaz309/todo-api/internal/handlers"
)

func Routes(app *app.App) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.WelcomePage)

	mux.HandleFunc("/categories", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		default:
			http.Error(w, "Wrong method", http.StatusMethodNotAllowed)
		case http.MethodGet:
			handlers.GetCategories(app)(w, r)
		case http.MethodPost:
			handlers.AddCategory(app)(w, r)
		}
	})

	mux.HandleFunc("/categories/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		default:
			http.Error(w, "Wrong method", http.StatusMethodNotAllowed)
		case http.MethodPut:
			handlers.UpdateCategory(app)(w, r)
		case http.MethodDelete:
			handlers.DeleteCategory(app)(w, r)
		}
	})

	return mux
}
