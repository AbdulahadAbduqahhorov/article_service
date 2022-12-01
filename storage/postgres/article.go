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
		res       []*article_service.Article
		u         sql.NullString
		d         sql.NullString
		tempTitle string
		tempBody  string
	)
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
			article.CreatedAt = u.String
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
		article       article_service.GetArticleByIdResponse
		tempFullname  *string
		tempTitle     string
		tempBody      string
		tempId        string
		tempCreatedAt string
		u             sql.NullString
		d             sql.NullString
		au            sql.NullString
		ad            sql.NullString
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
		&article.Id,
		&tempTitle,
		&tempBody,
		&article.CreatedAt,
		&au,
		&ad,
		&tempId,
		&tempFullname,
		&tempCreatedAt,
		&u,
		&d,
	)
	if err != nil {
		return nil, err
	}
	if au.Valid {
		article.UpdatedAt = au.String
	}
	if ad.Valid {
		article.DeletedAt = ad.String
	}
	if tempFullname != nil {
		article.Author = &article_service.GetArticleByIdResponse_Author{
			Id:       tempId,
			FullName: *tempFullname,
			CreatedAt: tempCreatedAt,
			UpdatedAt: u.String,
			DeletedAt: u.String,
		}
	}
	article.Content = &article_service.Content{Title: tempTitle, Body: tempBody}
	
	return &article, nil
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
