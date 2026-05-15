package models

import "time"

type User struct {
	ID        	uint      	`json:"id"`
	CreatedAt 	time.Time 	`json:"created_at"`
	UpdatedAt	time.Time 	`json:"updated_at"`

	Username	string    	`json:"username"`
	Password	string		`json:"password"`
}

type PaginatedUsersResponse struct {
	Cursor		uint		`json:"cursor"`
	Limit		int		`json:"limit"`
	Users		[]User		`json:"users"`
}
