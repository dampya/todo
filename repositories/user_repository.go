package repositories

import (
	"go/todo/models"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

// UserService.CreateUser
func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// UserService.GetUser
func (r *userRepository) GetOne(userID uint) (*models.User, error) {
	var user models.User

	if err := r.db.First(&user, userID).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// UserService.GetUsers
func (r *userRepository) GetAll(cursor uint, limit int) ([]models.User, uint, error) {
	var users []models.User

	query := r.db.
		Order("id ASC").
		Limit(limit)

	if cursor > 0 {
		query = query.Where("id > ?", cursor)
	}

	if err := query.Find(&users).Error; err != nil {
		return nil, 0, err
	}

	var nextCursor uint

	if len(users) > 0 {
		nextCursor = users[len(users)-1].ID
	}

	return users, nextCursor, nil
}

// UserService.UpdateUser
func (r *userRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

// UserService.DeleteUser
func (r *userRepository) Delete(userID uint) error {
	return r.db.Delete(&models.User{}, userID).Error
}

// UserService.DeleteUser
func (r *userRepository) DeleteUserTodos(userID uint) error {
	return r.db.
		Where("user_id = ?", userID).
		Delete(&models.Todo{}).Error
}
