package helpers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func GetUserID(r *http.Request) (uint, error) {
	userID := chi.URLParam(r, "userID")

	id, err := strconv.Atoi(userID)
	if err != nil {
		return 0, ErrInvalidID
	}

	return uint(id), nil
}

func GetTodoID(r *http.Request) (uint, error) {
	todoID := chi.URLParam(r, "todoID")

	id, err := strconv.Atoi(todoID)
	if err != nil {
		return 0, ErrInvalidID
	}

	return uint(id), nil
}
