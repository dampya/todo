package handlers

import (
	"net/http"
	"strconv"

	"go/todo/helpers"
	"go/todo/models"
	"go/todo/services"
)

type TodoHandler struct {
	todoService *services.TodoService
}

func NewTodoHandler(todoService *services.TodoService) *TodoHandler {
	return &TodoHandler{todoService: todoService}
}

func (h *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	userID, ok := helpers.GetUserID(w, r)
	if !ok {
		return
	}

	var todo models.Todo

	if err := helpers.DecodeJSON(r, &todo); err != nil {
		helpers.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}

	if err := h.todoService.CreateTodo(userID, &todo); err != nil {
		helpers.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	helpers.EncodeJSON(w, http.StatusCreated, todo)
}

func (h *TodoHandler) GetTodo(w http.ResponseWriter, r *http.Request) {
	userID, ok := helpers.GetUserID(w, r)
	if !ok {
		return
	}

	todoID, ok := helpers.GetTodoID(w, r)
	if !ok {
		return
	}

	todo, err := h.todoService.GetTodo(userID, todoID)
	if err != nil {
		helpers.WriteError(w, http.StatusForbidden, err.Error())
		return
	}

	helpers.EncodeJSON(w, http.StatusOK, todo)
}

func (h *TodoHandler) GetTodos(w http.ResponseWriter, r *http.Request) {
	userID, ok := helpers.GetUserID(w, r)
	if !ok {
		return
	}

	page := 1
	limit := 10

	if p := r.URL.Query().Get("page"); p != "" {
		if val, err := strconv.Atoi(p); err == nil && val > 0 {
			page = val
		}
	}

	if l := r.URL.Query().Get("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil && val > 0 {
			limit = val
		}
	}

	todos, total, err := h.todoService.GetTodos(userID, page, limit)
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, "failed to get todos")
		return
	}

	resp := models.PaginatedTodosResponse{
		Page:  page,
		Limit: limit,
		Total: total,
		Todos: todos,
	}

	helpers.EncodeJSON(w, http.StatusOK, resp)
}

func (h *TodoHandler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	userID, ok := helpers.GetUserID(w, r)
	if !ok {
		return
	}

	todoID, ok := helpers.GetTodoID(w, r)
	if !ok {
		return
	}

	var todo models.Todo

	if err := helpers.DecodeJSON(r, &todo); err != nil {
		helpers.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}

	todo.ID = todoID

	updatedTodo, err := h.todoService.UpdateTodo(userID, &todo)
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	helpers.EncodeJSON(w, http.StatusOK, updatedTodo)
}

func (h *TodoHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	userID, ok := helpers.GetUserID(w, r)
	if !ok {
		return
	}

	todoID, ok := helpers.GetTodoID(w, r)
	if !ok {
		return
	}

	if err := h.todoService.DeleteTodo(userID, todoID); err != nil {
		helpers.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	helpers.EncodeJSON(w, http.StatusNoContent, nil)
}
