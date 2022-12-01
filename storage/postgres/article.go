package postgres

import (
	"database/sql"
	"errors"

	"github.com/AbdulahadAbduqahhorov/gRPC/blogpost/article_service/genproto/article_service"
	"github.com/AbdulahadAbduqahhorov/gRPC/blogpost/article_service/storage"
	"github.com/jmoiron/sqlx"
)

type articleRepo struct {
	db *sqlx.DB
}

func NewArticleRepo(db *sqlx.DB) storage.ArticleRepoI {
	return articleRepo{
		db: db,
	}
}

func (stg articleRepo) CreateArticle(id string, req *article_service.CreateArticleRequest) error {

	_, err := stg.db.Exec(`INSERT INTO article 
	(
		id,
		title,
		body,
		author_id
	) VALUES (
		$1,
		$2,
		$3,
		$4
	)`,
		id,
		req.Content.Title,
		req.Content.Body,
		req.AuthorId,
	)
	if err != nil {
		return err
	}
	return nil

}

func (stg articleRepo) GetArticleList(req *article_service.GetArticleListRequest) (*article_service.GetArticleListResponse, error) {
	var (
		u         sql.NullString
		d         sql.NullString
		tempTitle string
		tempBody  string
	)
	res := make([]*article_service.Article, 0)
	rows, err := stg.db.Queryx(`SELECT 
		id,
		title,
		body,
		author_id,
		created_at,
		updated_at,
		deleted_at 
		FROM article
		WHERE title ILIKE '%' || $1 || '%' AND deleted_at IS NULL
		LIMIT $2
		OFFSET $3
	`,
		req.Search,
		int(req.Limit),
		int(req.Offset),
	)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var article article_service.Article
		err := rows.Scan(
			&article.Id,
			&tempTitle,
			&tempBody,
			&article.AuthorId,
			&article.CreatedAt,
			&u,
			&d,
		)
		if err != nil {
			return nil, err
		}
		if u.Valid {
			article.UpdatedAt = u.String
		}
		if d.Valid {
			article.DeletedAt = d.String
		}
		article.Content = &article_service.Content{Title: tempTitle, Body: tempBody}
		res = append(res, &article)

	}

	return &article_service.GetArticleListResponse{
		Articles: res,
	}, nil

}

func (stg articleRepo) GetArticleById(id string) (*article_service.GetArticleByIdResponse, error) {
	var (
		res article_service.GetArticleByIdResponse

		articleUpdated_at    sql.NullString
		articleDeleted_at    sql.NullString

		tempContentTitle     string
		tempContentBody      string

		tempAuthorId         string
		tempAuthorCreatedAt  string
		tempAuthorFullname   *string
		tempAuthorUpdatedAt  sql.NullString
		tempAuthorDeleteddAt sql.NullString
	)
	err := stg.db.QueryRow(`
	SELECT 
		ar.id,
		ar.title,
		ar.body,
		ar.created_at,
		ar.updated_at,
		ar.deleted_at,
		au.id,
		au.fullname,
		au.created_at,
		au.updated_at,
		au.deleted_at
	FROM article ar JOIN author au ON ar.author_id=au.id WHERE ar.id=$1 `, id).Scan(
		&res.Id,
		&tempContentTitle,
		&tempContentBody,
		&res.CreatedAt,
		&articleUpdated_at,
		&articleDeleted_at,
		&tempAuthorId,
		&tempAuthorFullname,
		&tempAuthorCreatedAt,
		&tempAuthorUpdatedAt,
		&tempAuthorDeleteddAt,
	)
	if err != nil {
		return nil, err
	}
	if articleUpdated_at.Valid {
		res.UpdatedAt = articleDeleted_at.String
	}
	if articleDeleted_at.Valid {
		res.DeletedAt = articleDeleted_at.String
	}
	var (
		fname string
		u string
		d string
	) 
	if tempAuthorFullname != nil {
		fname=*tempAuthorFullname
	}
	if tempAuthorUpdatedAt.Valid{
		u=tempAuthorUpdatedAt.String
	}

	if tempAuthorDeleteddAt.Valid{
		d=tempAuthorDeleteddAt.String
	}

	res.Content= &article_service.Content{
		Title: tempContentTitle,
		Body: tempContentBody,
	}

	res.Author= &article_service.GetArticleByIdResponse_Author{
		Id: tempAuthorId,
		FullName: fname,
		CreatedAt: tempAuthorCreatedAt,
		UpdatedAt: u,
		DeletedAt: d,
	}


	return &res, nil
}
	
func (stg articleRepo) UpdateArticle(req *article_service.UpdateArticleRequest) error {

	res, err := stg.db.NamedExec(`
	UPDATE  article SET 
		title=:t, 
		body=:b,
		updated_at=now() 
		WHERE id=:i AND deleted_at IS NULL 	`, map[string]interface{}{
		"t": req.Content.Title,
		"b": req.Content.Body,
		"i": req.Id,
	})
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n > 0 {
		return nil
	}
	return errors.New("article not found")
}

func (stg articleRepo) DeleteArticle(id string) error {
	res, err := stg.db.Exec(`UPDATE article SET deleted_at=now() WHERE id=$1 AND deleted_at IS NULL`, id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n > 0 {
		return nil
	}
	return errors.New("article not found")

}
