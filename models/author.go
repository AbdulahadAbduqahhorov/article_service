package models

import "time"

type Author struct {
	Id        string     `json:"id"`
	FullName  string     `json:"fullname"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type CreateAuthorModel struct {
	FullName string `json:"fullname" binding:"required"`
}

type UpdateAuthorModel struct {
	Id       string `json:"id"`
	FullName string `json:"fullname" binding:"required"`
}
