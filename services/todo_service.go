package services

import (
	"log/slog"

	"go/todo/helpers"
	"go/todo/models"
	"go/todo/repositories"
)

type TodoService struct {
	todoRepo repositories.TodoRepository
	userRepo repositories.UserRepository
}

func NewTodoService(
	todoRepo repositories.TodoRepository,
	userRepo repositories.UserRepository,
	) *TodoService {
	return &TodoService{
		todoRepo: todoRepo,
		userRepo: userRepo,
	}
}

// CreateTodo
func (s *TodoService) CreateTodo(userID uint, todo *models.Todo) error {
	if userID == 0 {
		return helpers.ErrInvalidID
	}

	_, err := s.userRepo.GetOne(userID)
	if err != nil {
		return helpers.ErrUserNotFound
	}

	if todo.Title == "" || todo.Description == "" {
		return helpers.ErrMissingFields
	}

	todo.UserID = userID

	if err := s.todoRepo.Create(todo); err != nil {
		return helpers.MapDBError(err)
	}

	slog.Info("todo created",
		"userID", userID,
		"todoID", todo.ID,
	)

	return nil
}

// GetTodo
func (s *TodoService) GetTodo(userID, todoID uint) (*models.Todo, error) {
	if userID == 0 || todoID == 0 {
		return nil, helpers.ErrInvalidID
	}

	todo, err := s.todoRepo.GetOne(todoID)
	if err != nil {
		return nil, helpers.ErrTodoNotFound
	}

	if todo.UserID != userID {
		return nil, helpers.ErrForbidden
	}

	slog.Info("todo fetched",
		"userID", userID,
		"todoID", todoID,
	)

	return todo, nil
}

// GetTodos
func (s *TodoService) GetTodos(userID uint, page, limit int) ([]models.Todo, int64, error) {
	if userID == 0 {
		return nil, 0, helpers.ErrInvalidID
	}

	if limit <= 0 || limit > 10 {
		limit = 10
	}

	todos, total, err := s.todoRepo.GetAll(userID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	slog.Info("todos fetched",
		"userID", userID,
		"count", len(todos),
	)

	return todos, total, nil
}

// UpdateTodo
func (s *TodoService) UpdateTodo(userID uint, todo *models.Todo) (*models.Todo, error) {
	if userID == 0 || todo.ID == 0 {
		return nil, helpers.ErrInvalidID
	}

	existing, err := s.todoRepo.GetOne(todo.ID)
	if err != nil {
		return nil, helpers.ErrTodoNotFound
	}

	if existing.UserID != userID {
		return nil, helpers.ErrForbidden
	}

	if todo.Title != "" {
		existing.Title = todo.Title
	}

	if todo.Description != "" {
		existing.Description = todo.Description
	}

	existing.Completed = todo.Completed

	if err := s.todoRepo.Update(existing); err != nil {
		return nil, helpers.MapDBError(err)
	}

	slog.Info("todo updated",
		"userID", userID,
		"todoID", todo.ID,
	)

	return existing, nil
}

// DeleteTodo
func (s *TodoService) DeleteTodo(userID, todoID uint) error {
	if userID == 0 || todoID == 0 {
		return helpers.ErrInvalidID
	}

	todo, err := s.todoRepo.GetOne(todoID)
	if err != nil {
		return helpers.ErrTodoNotFound
	}

	if todo.UserID != userID {
		return helpers.ErrForbidden
	}

	if err := s.todoRepo.Delete(todoID); err != nil {
		return helpers.MapDBError(err)
	}

	slog.Info("todo deleted",
		"userID", userID,
		"todoID", todoID,
	)

	return nil
}
