package repositories

import "go/todo/models"

type UserRepository interface {
	Create(user *models.User) error
	GetOne(userID uint) (*models.User, error)
	GetAll(cursor uint, limit int) ([]models.User, uint, error)
	Update(user *models.User) error
	Delete(userID uint) error
	DeleteUserTodos(userID uint) error
}

type TodoRepository interface {
	Create(todo *models.Todo) error
	GetOne(todoID uint) (*models.Todo, error)
	GetAll(userID uint, page int, limit int) ([]models.Todo, int64, error)
	Update(todo *models.Todo) error
	Delete(todoID uint) error
}
