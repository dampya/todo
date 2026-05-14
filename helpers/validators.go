package helpers

import (
	"net/http"
	"strconv"
	
	"github.com/go-chi/chi/v5"
)

func GetUserID(w http.ResponseWriter, r *http.Request) (uint, bool) {
	userID := chi.URLParam(r, "userID")

	id, err := strconv.Atoi(userID)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return 0, false
	}

	return uint(id), true
}

func GetTodoID(w http.ResponseWriter, r *http.Request) (uint, bool) {
	todoID := chi.URLParam(r, "todoID")
	
	id, err := strconv.Atoi(todoID)
	if err != nil {
		http.Error(w, "invalid todo id", http.StatusBadRequest)
		return 0, false
	}

	return uint(id), true
}
