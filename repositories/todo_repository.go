package repositories

import (
	"database/sql"
	"go/todo/models"
)

type todoRepository struct {
	db *sql.DB
}

func NewTodoRepository(db *sql.DB) TodoRepository {
	return &todoRepository{
		db: db,
	}
}

// Create Todo
func (r *todoRepository) Create(todo *models.Todo) error {
	query := `
		INSERT INTO todos
		(title, description, completed, user_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	return r.db.QueryRow(
		query,
		todo.Title,
		todo.Description,
		todo.Completed,
		todo.UserID,
	).Scan(
		&todo.ID,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)
}

// Get One Todo
func (r *todoRepository) GetOne(todoID uint) (*models.Todo, error) {
	query := `
		SELECT
			id,
			title,
			description,
			completed,
			user_id,
			created_at,
			updated_at
		FROM todos
		WHERE id = $1
	`

	var todo models.Todo

	err := r.db.QueryRow(query, todoID).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Description,
		&todo.Completed,
		&todo.UserID,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &todo, nil
}

// Get All Todos
func (r *todoRepository) GetAll(
	userID uint,
	page int,
	limit int,
) ([]models.Todo, int64, error) {

	offset := (page - 1) * limit

	countQuery := `
		SELECT COUNT(*)
		FROM todos
		WHERE user_id = $1
	`

	var total int64

	err := r.db.QueryRow(countQuery, userID).Scan(&total)

	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT
			id,
			title,
			description,
			completed,
			user_id,
			created_at,
			updated_at
		FROM todos
		WHERE user_id = $1
		ORDER BY id ASC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, userID, limit, offset)

	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()

	var todos []models.Todo

	for rows.Next() {
		var todo models.Todo

		err := rows.Scan(
			&todo.ID,
			&todo.Title,
			&todo.Description,
			&todo.Completed,
			&todo.UserID,
			&todo.CreatedAt,
			&todo.UpdatedAt,
		)

		if err != nil {
			return nil, 0, err
		}

		todos = append(todos, todo)
	}

	return todos, total, nil
}

// Update Todo
func (r *todoRepository) Update(todo *models.Todo) error {
	query := `
		UPDATE todos
		SET
			title = $1,
			description = $2,
			completed = $3,
			updated_at = NOW()
		WHERE id = $4
	`

	_, err := r.db.Exec(
		query,
		todo.Title,
		todo.Description,
		todo.Completed,
		todo.ID,
	)

	return err
}

// Delete Todo
func (r *todoRepository) Delete(todoID uint) error {
	query := `DELETE FROM todos WHERE id = $1`

	_, err := r.db.Exec(query, todoID)

	return err
}
