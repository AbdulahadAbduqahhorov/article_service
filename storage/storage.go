package storage

import (
	"github.com/AbdulahadAbduqahhorov/gin/Article/genproto/author_service"
	"github.com/AbdulahadAbduqahhorov/gin/Article/models"
)

type StorageI interface {
	Article() ArticleRepoI
	Author() AuthorRepoI
}

type ArticleRepoI interface {
	CreateArticle(id string, article models.CreateArticleModel) error
	GetArticle(limit, offset int, search string) ([]models.Article, error)
	GetArticleById(id string) (models.GetArticleByIdModel, error)
	UpdateArticle(article models.UpdateArticleModel) error
	DeleteArticle(id string) error
}

type AuthorRepoI interface {
	CreateAuthor(author_service.CreateAuthorRequest) (author_service.CreateAuthorResponse,error)
	GetAuthor(author_service.GetAuthorRequest) (author_service.GetAuthorResponse, error)
	GetAuthorById(author_service.GetAuthorByIdResponse) (author_service.Author, error)
	UpdateAuthor(author models.UpdateAuthorModel) error
	DeleteAuthor(id string) error
}
