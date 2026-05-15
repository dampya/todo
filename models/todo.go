package models

import "time"

type Todo struct {
	ID		uint		`json:"id"`
	CreatedAt	time.Time	`json:"created_at"`
	UpdatedAt	time.Time	`json:"updated_at"`

	Title		string		`json:"title"`
	Description	string		`json:"description"`
	Completed	bool		`json:"completed"`
	UserID		uint		`json:"user_id"`
}

type PaginatedTodosResponse struct {
	Page		int		`json:"page"`
	Limit		int		`json:"limit"`
	Total		int64		`json:"total"`
	Todos		[]Todo		`json:"todos"`
}
