package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"go/todo/database"
	"go/todo/handlers"
	"go/todo/repositories"
	"go/todo/routes"
	"go/todo/services"
	"go/todo/telemetry"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	file, err := os.OpenFile(
		"/home/om/projects/todo/telemetry/log/app.log",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0666,
	)
	if err != nil {
		panic(err)
	}

	logger := slog.New(
		slog.NewJSONHandler(file, nil),
	)
	slog.SetDefault(logger)

	db := database.ConnectDB()

	userRepo := repositories.NewUserRepository(db)
	todoRepo := repositories.NewTodoRepository(db)

	userService := services.NewUserService(userRepo)
	todoService := services.NewTodoService(todoRepo, userRepo)

	userHandler := handlers.NewUserHandler(userService)
	todoHandler := handlers.NewTodoHandler(todoService)

	r := chi.NewRouter()

	shutdown := telemetry.Init()
	defer shutdown()

	r.Use(telemetry.Middleware)

	routes.SetupRoutes(r, userHandler, todoHandler)

	fmt.Println("Server running on http://127.0.0.1:8080")
	if err := http.ListenAndServe("127.0.0.1:8080", r); err != nil {
		slog.Error("server failed", "error", err)
	}
}
