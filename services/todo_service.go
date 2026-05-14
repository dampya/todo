package services

import (
	"errors"
	"log/slog"

	"go/todo/models"
	"go/todo/repositories"
)

type TodoService struct {
	todoRepo repositories.TodoRepository
	userRepo repositories.UserRepository
}

func NewTodoService(todoRepo repositories.TodoRepository, userRepo repositories.UserRepository) *TodoService {
	return &TodoService{
		todoRepo: todoRepo,
		userRepo: userRepo,
	}
}

// TodoHandler.CreateTodo
func (s *TodoService) CreateTodo(userID uint, todo *models.Todo) error {
	if userID == 0 {
		return errors.New("invalid id")
	}

	_, err := s.userRepo.GetOne(userID)
	if err != nil {
		return errors.New("user not found")
	}

	if todo.Title == "" || todo.Description == "" {
		return errors.New("missing fields")
	}

	todo.UserID = userID

	if err := s.todoRepo.Create(todo); err != nil {
		return err
	}

	slog.Info("todo created", "userID", userID, "todoID", todo.ID)

	return nil
}

// TodoHandler.GetTodo
func (s *TodoService) GetTodo(userID, todoID uint) (*models.Todo, error) {
	if userID == 0 || todoID == 0 {
		return nil, errors.New("invalid id")
	}

	todo, err := s.todoRepo.GetOne(todoID)
	if err != nil {
		return nil, err
	}

	if todo.UserID != userID {
		return nil, errors.New("forbidden")
	}

	slog.Info("todo fetched", "userID", userID, "todoID", todoID)

	return todo, nil
}

// TodoHandler.GetTodos
func (s *TodoService) GetTodos(userID uint, page, limit int) ([]models.Todo, int64, error) {
	if userID == 0 {
		return nil, 0, errors.New("invalid id")
	}

	if limit <= 0 || limit > 10 {
		limit = 10
	}

	todos, total, err := s.todoRepo.GetAll(userID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	slog.Info("todos fetched", "userID", userID, "count", len(todos))

	return todos, total, nil
}

// TodoHandler.UpdateTodo
func (s *TodoService) UpdateTodo(userID uint, todo *models.Todo) (*models.Todo, error) {
	if userID == 0 || todo.ID == 0 {
		return nil, errors.New("invalid id")
	}

	existing, err := s.todoRepo.GetOne(todo.ID)
	if err != nil {
		return nil, err
	}

	if existing.UserID != userID {
		return nil, errors.New("forbidden")
	}

	if todo.Title != "" {
		existing.Title = todo.Title
	}

	if todo.Description != "" {
		existing.Description = todo.Description
	}

	existing.Completed = todo.Completed

	if err := s.todoRepo.Update(existing); err != nil {
		return nil, err
	}

	slog.Info("todo updated", "userID", userID, "todoID", todo.ID)

	return existing, nil
}

// TodoHandler.DeleteTodo
func (s *TodoService) DeleteTodo(userID, todoID uint) error {
	if userID == 0 || todoID == 0 {
		return errors.New("invalid id")
	}

	todo, err := s.todoRepo.GetOne(todoID)
	if err != nil {
		return err
	}

	if todo.UserID != userID {
		return errors.New("forbidden")
	}

	if err := s.todoRepo.Delete(todoID); err != nil {
		return err
	}

	slog.Info("todo deleted", "userID", userID, "todoID", todoID)

	return nil
}
