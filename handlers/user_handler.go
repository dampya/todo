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

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User

	if err := helpers.DecodeJSON(r, &user); err != nil {
		helpers.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}

	if err := h.userService.CreateUser(&user); err != nil {
		helpers.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	helpers.EncodeJSON(w, http.StatusCreated, user)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := helpers.GetUserID(w, r)
	if !ok {
		return
	}

	user, err := h.userService.GetUser(userID)
	if err != nil {
		helpers.WriteError(w, http.StatusNotFound, err.Error())
		return
	}

	helpers.EncodeJSON(w, http.StatusOK, user)
}

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
		helpers.WriteError(w, http.StatusInternalServerError, "failed to get users")
		return
	}

	resp := models.PaginatedUsersResponse{
		Cursor: nextCursor,
		Limit:  limit,
		Users:  users,
	}

	helpers.EncodeJSON(w, http.StatusOK, resp)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := helpers.GetUserID(w, r)
	if !ok {
		return
	}

	var user models.User

	if err := helpers.DecodeJSON(r, &user); err != nil {
		helpers.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}

	user.ID = userID

	updatedUser, err := h.userService.UpdateUser(userID, &user)
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	helpers.EncodeJSON(w, http.StatusOK, updatedUser)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := helpers.GetUserID(w, r)
	if !ok {
		return
	}

	if err := h.userService.DeleteUser(userID); err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.EncodeJSON(w, http.StatusNoContent, nil)
}
