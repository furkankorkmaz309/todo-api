package routes

import (
	"net/http"

	"github.com/furkankorkmaz309/todo-api/internal/app"
	"github.com/furkankorkmaz309/todo-api/internal/handlers"
	"github.com/go-chi/chi"
)

func Routes(app *app.App) http.Handler {
	r := chi.NewRouter()

	r.Use(
		handlers.RecoverPanic(app),
		handlers.LimitRequest(app),
		handlers.LogRequest(app))

	r.Get("/", handlers.WelcomePage)

	r.Route("/categories", func(r chi.Router) {
		r.Get("/", handlers.GetCategories(app))
		r.Post("/", handlers.AddCategory(app))
		r.Patch("/{id}", handlers.PatchCategory(app))
		r.Delete("/{id}", handlers.DeleteCategory(app))
	})

	r.Route("/todos", func(r chi.Router) {
		r.Get("/", handlers.GetTodos(app))
		r.Post("/", handlers.CreateTodo(app))
		r.Get("/{id}", handlers.GetTodo(app))
		r.Patch("/{id}", handlers.PatchTodo(app))
		r.Delete("/{id}", handlers.DeleteTodo(app))
	})

	return r
}
