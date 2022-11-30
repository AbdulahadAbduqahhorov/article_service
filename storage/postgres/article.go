package postgres

import (
	"errors"

	"github.com/AbdulahadAbduqahhorov/gRPC/blogpost/article_service/models"
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

func (stg articleRepo) CreateArticle(id string, article models.CreateArticleModel) error {

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
		article.Title,
		article.Body,
		article.AuthorId,
	)
	if err != nil {
		return err
	}
	return nil

}

func (stg articleRepo) GetArticle(limit, offset int, search string) ([]models.Article, error) {
	var res []models.Article

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
		search,
		limit,
		offset,
	)
	if err != nil {
		return res, err
	}

	for rows.Next() {
		var article models.Article
		err := rows.Scan(
			&article.Id,
			&article.Title,
			&article.Body,
			&article.AuthorId,
			&article.CreatedAt,
			&article.UpdatedAt,
			&article.DeletedAt,
		)
		if err != nil {
			return res, err
		}

		res = append(res, article)

	}

	return res, err

}

func (stg articleRepo) GetArticleById(id string) (models.GetArticleByIdModel, error) {
	var article models.GetArticleByIdModel
	var tempFullname *string
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
		&article.Title,
		&article.Body,
		&article.CreatedAt,
		&article.UpdatedAt,
		&article.DeletedAt,
		&article.Author.Id,
		&article.Author.FullName,
		&article.Author.CreatedAt,
		&article.Author.UpdatedAt,
		&article.Author.DeletedAt,
	)
	if err != nil {
		return article, err
	}

	if tempFullname != nil {
		article.Author.FullName = *tempFullname
	}

	return article, nil
}

func (stg articleRepo) UpdateArticle(article models.UpdateArticleModel) error {

	res, err := stg.db.NamedExec(`
	UPDATE  article SET 
		title=:t, 
		body=:b,
		updated_at=now() 
		WHERE id=:i AND deleted_at IS NULL 	`, map[string]interface{}{
		"t": article.Title,
		"b": article.Body,
		"i": article.Id,
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
