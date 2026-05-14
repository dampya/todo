package models

import "gorm.io/gorm"

type Todo struct {
	gorm.Model
	Title		string		`json:"title"`
	Description	string		`json:"description"`
	Completed	bool		`json:"completed"`
	UserID		uint		`gorm:"foreignKey:UserID"`
}

type PaginatedTodosResponse struct {
	Page		int		`json:"page"`
	Limit		int		`json:"limit"`
	Total		int64		`json:"total"`
	Todos		[]Todo		`json:"todos"`
}
