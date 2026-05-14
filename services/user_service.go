package services

import (
	"errors"
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

// UserHandler.CreateUser
func (s *UserService) CreateUser(user *models.User) error {
	if user.Username == "" || user.Password == "" {
		return errors.New("missing fields")
	}

	if err := s.userRepo.Create(user); err != nil {
		if helpers.IsDuplicateError(err) {
			return errors.New("username already exists")
		}
		return err
	}

	slog.Info("user created", "userID", user.ID, "username", user.Username)

	return nil
}

// UserHandler.GetUser
func (s *UserService) GetUser(userID uint) (*models.User, error) {
	if userID == 0 {
		return nil, errors.New("invalid id")
	}

	user, err := s.userRepo.GetOne(userID)
	if err != nil {
		return nil, err
	}

	slog.Info("user fetched", "userID", userID)

	return user, nil
}

// UserHandler.GetUsers
func (s *UserService) GetUsers(cursor uint, limit int) ([]models.User, uint, error) {
	if limit <= 0 || limit > 10 {
		limit = 10
	}

	users, nextCursor, err := s.userRepo.GetAll(cursor, limit)
	if err != nil {
		return nil, 0, err
	}

	slog.Info("users fetched", "count", len(users), "nextCursor", nextCursor)

	return users, nextCursor, nil
}

// UserHandler.UpdateUser
func (s *UserService) UpdateUser(userID uint, user *models.User) (*models.User, error) {
	if userID == 0 {
		return nil, errors.New("invalid id")
	}

	if user.Username == "" || user.Password == "" {
		return nil, errors.New("fields cannot be empty")
	}

	existing, err := s.userRepo.GetOne(userID)
	if err != nil {
		return nil, err
	}

	existing.Username = user.Username
	existing.Password = user.Password

	if err := s.userRepo.Update(existing); err != nil {
		return nil, err
	}

	slog.Info("user updated", "userID", userID)

	return existing, nil
}

// UserHandler.DeleteUser
func (s *UserService) DeleteUser(userID uint) error {
	if userID == 0 {
		return errors.New("invalid id")
	}

	if err := s.userRepo.DeleteUserTodos(userID); err != nil {
		return err
	}

	if err := s.userRepo.Delete(userID); err != nil {
		return err
	}

	slog.Info("user deleted", "userID", userID)

	return nil
}
