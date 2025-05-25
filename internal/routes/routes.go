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
		case http.MethodPatch:
			handlers.PatchCategory(app)(w, r)
		case http.MethodDelete:
			handlers.DeleteCategory(app)(w, r)
		}
	})

	mux.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		default:
			http.Error(w, "Wrong method", http.StatusMethodNotAllowed)
		case http.MethodGet:
			handlers.GetTodos(app)(w, r)
		case http.MethodPost:
			handlers.CreateTodo(app)(w, r)
		}
	})

	mux.HandleFunc("/todos/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		default:
			http.Error(w, "Wrong method", http.StatusMethodNotAllowed)
		case http.MethodGet:
			handlers.GetTodo(app)(w, r)
		case http.MethodPatch:
			handlers.PatchTodo(app)(w, r)
		case http.MethodDelete:
			handlers.DeleteTodo(app)(w, r)
		}
	})

	return mux
}
