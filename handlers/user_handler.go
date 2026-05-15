package handlers

import (
	"net/http"
	"strconv"

	"go/todo/helpers"
	"go/todo/models"
	"go/todo/services"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// POST /user
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User

	if err := helpers.DecodeJSON(r, &user); err != nil {
		helpers.HandleError(w, err)
		return
	}

	if err := h.userService.CreateUser(&user); err != nil {
		helpers.HandleError(w, err)
		return
	}

	helpers.EncodeJSON(w, http.StatusCreated, user)
}

// GET /user/{userID}
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userID, err := helpers.GetUserID(r)
	if err != nil {
		helpers.HandleError(w, err)
		return
	}

	user, err := h.userService.GetUser(userID)
	if err != nil {
		helpers.HandleError(w, err)
		return
	}

	helpers.EncodeJSON(w, http.StatusOK, user)
}

// GET /user
func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	cursorStr := r.URL.Query().Get("cursor")
	limitStr := r.URL.Query().Get("limit")

	var cursor uint
	var limit = 10

	if cursorStr != "" {
		c, _ := strconv.ParseUint(cursorStr, 10, 64)
		cursor = uint(c)
	}

	if limitStr != "" {
		l, _ := strconv.Atoi(limitStr)
		if l > 0 {
			limit = l
		}
	}

	users, nextCursor, err := h.userService.GetUsers(cursor, limit)
	if err != nil {
		helpers.HandleError(w, err)
		return
	}

	resp := models.PaginatedUsersResponse{
		Cursor: nextCursor,
		Limit:  limit,
		Users:  users,
	}

	helpers.EncodeJSON(w, http.StatusOK, resp)
}

// PUT /user/{userID}
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID, err := helpers.GetUserID(r)
	if err != nil {
		helpers.HandleError(w, err)
		return
	}

	var user models.User

	if err := helpers.DecodeJSON(r, &user); err != nil {
		helpers.HandleError(w, err)
		return
	}

	user.ID = userID

	updatedUser, err := h.userService.UpdateUser(userID, &user)
	if err != nil {
		helpers.HandleError(w, err)
		return
	}

	helpers.EncodeJSON(w, http.StatusOK, updatedUser)
}

// DELETE /user/{userID}
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID, err := helpers.GetUserID(r)
	if err != nil {
		helpers.HandleError(w, err)
		return
	}

	if err := h.userService.DeleteUser(userID); err != nil {
		helpers.HandleError(w, err)
		return
	}

	helpers.EncodeJSON(w, http.StatusNoContent, nil)
}
