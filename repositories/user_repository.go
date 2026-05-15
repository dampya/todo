package repositories

import (
	"database/sql"
	"go/todo/models"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

// Create User
func (r *userRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users
		(username, password, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	return r.db.QueryRow(
		query,
		user.Username,
		user.Password,
	).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
}

// Get One User
func (r *userRepository) GetOne(userID uint) (*models.User, error) {
	query := `
		SELECT
			id,
			username,
			password,
			created_at,
			updated_at
		FROM users
		WHERE id = $1
	`

	var user models.User

	err := r.db.QueryRow(query, userID).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Get All Users (Cursor Pagination)
func (r *userRepository) GetAll(
	cursor uint,
	limit int,) ([]models.User, uint, error) {

	var (
		rows *sql.Rows
		err  error
	)

	if cursor > 0 {
		query := `
			SELECT
				id,
				username,
				password,
				created_at,
				updated_at
			FROM users
			WHERE id > $1
			ORDER BY id ASC
			LIMIT $2
		`

		rows, err = r.db.Query(query, cursor, limit)

	} else {
		query := `
			SELECT
				id,
				username,
				password,
				created_at,
				updated_at
			FROM users
			ORDER BY id ASC
			LIMIT $1
		`

		rows, err = r.db.Query(query, limit)
	}

	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User

		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Password,
			&user.CreatedAt,
			&user.UpdatedAt,
		)

		if err != nil {
			return nil, 0, err
		}

		users = append(users, user)
	}

	var nextCursor uint

	if len(users) > 0 {
		nextCursor = users[len(users)-1].ID
	}

	return users, nextCursor, nil
}

// Update User
func (r *userRepository) Update(user *models.User) error {
	query := `
		UPDATE users
		SET
			username = $1,
			password = $2,
			updated_at = NOW()
		WHERE id = $3
	`

	_, err := r.db.Exec(
		query,
		user.Username,
		user.Password,
		user.ID,
	)

	return err
}

// Delete User
func (r *userRepository) Delete(userID uint) error {
	query := `DELETE FROM users WHERE id = $1`

	_, err := r.db.Exec(query, userID)

	return err
}

// Delete User Todos
func (r *userRepository) DeleteUserTodos(userID uint) error {
	query := `DELETE FROM todos WHERE user_id = $1`

	_, err := r.db.Exec(query, userID)

	return err
}
