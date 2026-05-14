package repositories

import (
	"go/todo/models"

	"gorm.io/gorm"
)

type todoRepository struct {
	db *gorm.DB
}

func NewTodoRepository(db *gorm.DB) TodoRepository {
	return &todoRepository{
		db: db,
	}
}

// TodoService.CreateTodo
func (r *todoRepository) Create(todo *models.Todo) error {
	return r.db.Create(todo).Error
}

// TodoService.GetTodo
func (r *todoRepository) GetOne(todoID uint) (*models.Todo, error) {
	var todo models.Todo

	if err := r.db.First(&todo, todoID).Error; err != nil {
		return nil, err
	}

	return &todo, nil
}

// TodoService.GetTodos
func (r *todoRepository) GetAll(userID uint, page int, limit int) ([]models.Todo, int64, error) {
	var todos []models.Todo
	var total int64

	offset := (page - 1) * limit

	if err := r.db.
		Model(&models.Todo{}).
		Where("user_id = ?", userID).
		Count(&total).Error; err != nil {

		return nil, 0, err
	}

	if err := r.db.
		Where("user_id = ?", userID).
		Limit(limit).
		Offset(offset).
		Order("id ASC").
		Find(&todos).Error; err != nil {

		return nil, 0, err
	}

	return todos, total, nil
}

// TodoService.UpdateTodo
func (r *todoRepository) Update(todo *models.Todo) error {
	return r.db.Save(todo).Error
}

// TodoService.DeleteTodo
func (r *todoRepository) Delete(todoID uint) error {
	return r.db.Delete(&models.Todo{}, todoID).Error
}
