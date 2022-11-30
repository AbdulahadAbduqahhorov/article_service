package models

import "time"

type Content struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type Article struct {
	Id string `json:"id"`
	Content
	AuthorId  string     `json:"author_id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"-"`
}

type CreateArticleModel struct {
	Content
	AuthorId string `json:"author_id"`
}

type UpdateArticleModel struct {
	Id string `json:"id"`
	Content
}

type GetArticleByIdModel struct {
	Id string `json:"id"`
	Content
	Author    Author     `json:"author"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
