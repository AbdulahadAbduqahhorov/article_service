package storage

import (
	"github.com/AbdulahadAbduqahhorov/gRPC/blogpost/article_service/genproto/article_service"
	"github.com/AbdulahadAbduqahhorov/gRPC/blogpost/article_service/genproto/author_service"
)

type StorageI interface {
	Article() ArticleRepoI
	Author() AuthorRepoI
}

type ArticleRepoI interface {
	CreateArticle(id string, req *article_service.CreateArticleRequest) error
	GetArticleList(req *article_service.GetArticleListRequest) (*article_service.GetArticleListResponse, error)
	GetArticleById(id string) (*article_service.GetArticleByIdResponse, error)
	UpdateArticle(req *article_service.UpdateArticleRequest) error
	DeleteArticle(id string) error
}

type AuthorRepoI interface {
	CreateAuthor(req *author_service.CreateAuthorRequest) (string, error)
	GetAuthor(req *author_service.GetAuthorRequest) (*author_service.GetAuthorResponse, error)
	GetAuthorById(id string) (*author_service.Author, error)
	UpdateAuthor(req *author_service.UpdateAuthorRequest)error
	DeleteAuthor(id string) error
}



