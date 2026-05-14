package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username	string		`json:"username" gorm:"uniqueIndex:idx_users_username"`
	Password	string		`json:"password"`
}

type PaginatedUsersResponse struct {
    Cursor 		uint   		`json:"cursor"`
    Limit  		int    		`json:"limit"`
    Users  		[]User 		`json:"users"`
}
