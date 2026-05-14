package routes

import (
	"go/todo/handlers"

	"github.com/go-chi/chi/v5"
)

func SetupRoutes(
	r chi.Router,
	userHandler *handlers.UserHandler,
	todoHandler *handlers.TodoHandler) {
	r.Route("/user", func(r chi.Router) {
		r.Get("/", userHandler.GetUsers)
		r.Post("/", userHandler.CreateUser)

		r.Route("/{userID}", func(r chi.Router) {
			r.Get("/", userHandler.GetUser)
			r.Put("/", userHandler.UpdateUser)
			r.Delete("/", userHandler.DeleteUser)

			r.Route("/todo", func(r chi.Router) {
				r.Get("/", todoHandler.GetTodos)
				r.Post("/", todoHandler.CreateTodo)

				r.Route("/{todoID}", func(r chi.Router) {
					r.Get("/", todoHandler.GetTodo)
					r.Put("/", todoHandler.UpdateTodo)
					r.Delete("/", todoHandler.DeleteTodo)
				})
			})
		})
	})
}
