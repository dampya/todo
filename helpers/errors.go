package helpers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/lib/pq"
)

var (
	ErrInvalidJSON    = errors.New("invalid json")
	ErrInvalidID      = errors.New("invalid id")
	ErrMissingFields  = errors.New("missing fields")
	ErrEmptyFields    = errors.New("fields cannot be empty")
	ErrUserNotFound   = errors.New("user not found")
	ErrTodoNotFound   = errors.New("todo not found")
	ErrUsernameExists = errors.New("username already exists")
	ErrForbidden      = errors.New("forbidden")
)

func writeError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(map[string]string{
		"error": msg,
	})
}

func HandleError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	switch {
	case errors.Is(err, ErrInvalidJSON):
		writeError(w, http.StatusBadRequest, err.Error())

	case errors.Is(err, ErrInvalidID):
		writeError(w, http.StatusBadRequest, err.Error())

	case errors.Is(err, ErrMissingFields):
		writeError(w, http.StatusBadRequest, err.Error())

	case errors.Is(err, ErrEmptyFields):
		writeError(w, http.StatusBadRequest, err.Error())

	case errors.Is(err, ErrUserNotFound):
		writeError(w, http.StatusNotFound, err.Error())

	case errors.Is(err, ErrTodoNotFound):
		writeError(w, http.StatusNotFound, err.Error())

	case errors.Is(err, ErrUsernameExists):
		writeError(w, http.StatusConflict, err.Error())

	case errors.Is(err, ErrForbidden):
		writeError(w, http.StatusForbidden, err.Error())

	default:
		writeError(w, http.StatusInternalServerError, "internal server error")
	}
}

func MapDBError(err error) error {
	if err == nil {
		return nil
	}

	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		switch pqErr.Code {
		case "23505": //psql unique violation code
			return ErrUsernameExists
		}
	}

	return err
}
