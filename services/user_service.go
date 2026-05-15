package services

import (
	"log/slog"

	"go/todo/helpers"
	"go/todo/models"
	"go/todo/repositories"
)

type UserService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

// CreateUser
func (s *UserService) CreateUser(user *models.User) error {
	if user.Username == "" || user.Password == "" {
		return helpers.ErrMissingFields
	}

	if err := s.userRepo.Create(user); err != nil {
		return helpers.MapDBError(err)
	}

	slog.Info("user created",
		"userID", user.ID,
		"username", user.Username,
	)

	return nil
}

// GetUser
func (s *UserService) GetUser(userID uint) (*models.User, error) {
	if userID == 0 {
		return nil, helpers.ErrInvalidID
	}

	user, err := s.userRepo.GetOne(userID)
	if err != nil {
		return nil, helpers.ErrUserNotFound
	}

	slog.Info("user fetched", "userID", userID)

	return user, nil
}

// GetUsers
func (s *UserService) GetUsers(cursor uint, limit int) ([]models.User, uint, error) {
	if limit <= 0 || limit > 10 {
		limit = 10
	}

	users, nextCursor, err := s.userRepo.GetAll(cursor, limit)
	if err != nil {
		return nil, 0, err
	}

	slog.Info("users fetched",
		"count", len(users),
		"nextCursor", nextCursor,
	)

	return users, nextCursor, nil
}

// UpdateUser
func (s *UserService) UpdateUser(userID uint, user *models.User) (*models.User, error) {
	if userID == 0 {
		return nil, helpers.ErrInvalidID
	}

	if user.Username == "" || user.Password == "" {
		return nil, helpers.ErrEmptyFields
	}

	existing, err := s.userRepo.GetOne(userID)
	if err != nil {
		return nil, helpers.ErrUserNotFound
	}

	existing.Username = user.Username
	existing.Password = user.Password

	if err := s.userRepo.Update(existing); err != nil {
		return nil, helpers.DBError(err)
	}

	slog.Info("user updated", "userID", userID)

	return existing, nil
}

// DeleteUser
func (s *UserService) DeleteUser(userID uint) error {
	if userID == 0 {
		return helpers.ErrInvalidID
	}

	// ensure user exists
	_, err := s.userRepo.GetOne(userID)
	if err != nil {
		return helpers.ErrUserNotFound
	}

	if err := s.userRepo.DeleteUserTodos(userID); err != nil {
		return helpers.DBError(err)
	}

	if err := s.userRepo.Delete(userID); err != nil {
		return helpers.DBError(err)
	}

	slog.Info("user deleted", "userID", userID)

	return nil
}
