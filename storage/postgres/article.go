package postgres

import (
	"database/sql"
	"errors"
	"time"

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

func (a articleRepo) CreateArticle(id string, req *article_service.CreateArticleRequest) error {
	if req.Content == nil {
		req.Content = &article_service.Content{}
	}
	_, err := a.db.Exec(`INSERT INTO article 
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

func (a articleRepo) GetArticleList(req *article_service.GetArticleListRequest) (*article_service.GetArticleListResponse, error) {

	res := &article_service.GetArticleListResponse{
		Articles: make([]*article_service.Article, 0),
	}
	rows, err := a.db.Queryx(`SELECT 
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
		article := &article_service.Article{
			Content: &article_service.Content{},
		}
		var updatedAt, deletedAt sql.NullString

		err := rows.Scan(
			&article.Id,
			&article.Content.Title,
			&article.Content.Body,
			&article.AuthorId,
			&article.CreatedAt,
			&updatedAt,
			&deletedAt,
		)
		if err != nil {
			return nil, err
		}
		if updatedAt.Valid {
			article.UpdatedAt = updatedAt.String
		}
		if deletedAt.Valid {
			article.DeletedAt = deletedAt.String
		}
		res.Articles = append(res.Articles, article)

	}

	return res, nil

}

func (a articleRepo) GetArticleById(id string) (*article_service.GetArticleByIdResponse, error) {
	res := &article_service.GetArticleByIdResponse{
		Content: &article_service.Content{},
		Author:  &article_service.GetArticleByIdResponse_Author{},
	}
	var (
		updatedAt sql.NullString
		deletedAt *time.Time

		authorUpdatedAt sql.NullString
		authorDeletedAt sql.NullString
	)
	err := a.db.QueryRow(`
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
		&res.Content.Title,
		&res.Content.Body,
		&res.CreatedAt,
		&updatedAt,
		&deletedAt,
		&res.Author.Id,
		&res.Author.FullName,
		&res.Author.CreatedAt,
		&authorUpdatedAt,
		&authorDeletedAt,
	)
	if err != nil {
		return nil, err
	}
	if deletedAt != nil {
		return nil, errors.New("article not found")
	}
	if updatedAt.Valid {
		res.UpdatedAt = updatedAt.String
	}

	if authorUpdatedAt.Valid {
		res.Author.UpdatedAt = authorUpdatedAt.String
	}
	if authorDeletedAt.Valid {
		res.Author.DeletedAt = authorDeletedAt.String
	}

	return res, nil
}

func (a articleRepo) UpdateArticle(req *article_service.UpdateArticleRequest) error {
	if req.Content == nil {
		req.Content = &article_service.Content{}
	}
	res, err := a.db.NamedExec(`
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

func (a articleRepo) DeleteArticle(id string) error {
	res, err := a.db.Exec(`UPDATE article SET deleted_at=now() WHERE id=$1 AND deleted_at IS NULL`, id)
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
