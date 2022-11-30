package postgres

import (
	"github.com/AbdulahadAbduqahhorov/gRPC/blogpost/article_service/storage"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Postgres struct {
	db      *sqlx.DB
	article storage.ArticleRepoI
	author  storage.AuthorRepoI
}

func NewPostgres(config string) (storage.StorageI, error) {
	tempDb, err := sqlx.Connect("postgres", config)
	if err != nil {
		return nil, err
	}
	return &Postgres{db: tempDb}, nil
}

func (p *Postgres) Article() storage.ArticleRepoI {
	if p.article == nil {
		p.article = NewArticleRepo(p.db)
	}
	return p.article
}

func (p *Postgres) Author() storage.AuthorRepoI {
	if p.author == nil {
		p.author = NewAuthorRepo(p.db)
	}
	return p.author
}
