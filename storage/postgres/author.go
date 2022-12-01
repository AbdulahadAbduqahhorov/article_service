package postgres

import (
	"database/sql"
	"errors"

	"github.com/AbdulahadAbduqahhorov/gRPC/blogpost/article_service/genproto/author_service"
	"github.com/AbdulahadAbduqahhorov/gRPC/blogpost/article_service/storage"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type authorRepo struct {
	db *sqlx.DB
}

func NewAuthorRepo(db *sqlx.DB) storage.AuthorRepoI {
	return authorRepo{
		db: db,
	}
}
func (stg authorRepo) CreateAuthor(req *author_service.CreateAuthorRequest) (string, error) {
	id := uuid.New().String()
	_, err := stg.db.Exec(`INSERT INTO 
		author (
			id,
			fullname
			) 
		VALUES (
			$1, 
			$2
			)`,
		id,
		req.FullName,
	)
	if err != nil {
		return "", err
	}
	return id, nil

}

func (stg authorRepo) GetAuthor(req *author_service.GetAuthorRequest) (*author_service.GetAuthorResponse, error) {
	var (
		authors      []*author_service.Author
		tempFullname *string
		u            sql.NullString
		d            sql.NullString
	)
	rows, err := stg.db.Queryx(`SELECT 
		id,
		fullname,
		created_at,
		updated_at,
		deleted_at 
		FROM author
		WHERE (fullname ILIKE '%' || $1 || '%') AND deleted_at IS NULL
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
		var author author_service.Author
		err := rows.Scan(
			&author.Id,
			&tempFullname,
			&author.CreatedAt,
			&u,
			&d,
		)
		if err != nil {
			return nil, err
		}
		if tempFullname != nil {
			author.FullName = *tempFullname
		}
		if u.Valid {
			author.UpdatedAt = u.String
		}
		if d.Valid {
			author.DeletedAt = d.String
		}
		authors = append(authors, &author)

	}

	return &author_service.GetAuthorResponse{Authors: authors}, err

}

func (stg authorRepo) GetAuthorById(id string) (*author_service.Author, error) {
	var (
		res          author_service.Author
		tempFullname *string

		u sql.NullString
		d sql.NullString
	)
	err := stg.db.QueryRow(`
	SELECT 
		id,
		fullname,
		created_at,
		updated_at,
		deleted_at
	FROM author  
	WHERE id=$1 AND deleted_at is NULL`, id).Scan(
		&res.Id,
		&tempFullname,
		&res.CreatedAt,
		&u,
		&d,
	)
	if u.Valid {
		res.UpdatedAt = u.String
	}
	if d.Valid {
		res.DeletedAt = d.String
	}
	if err != nil {
		return nil, err
	}
	if tempFullname != nil {
		res.FullName = *tempFullname
	}
	return &res, nil
}

func (stg authorRepo) UpdateAuthor(req *author_service.UpdateAuthorRequest) error {
	res, err := stg.db.NamedExec(`
	UPDATE  author SET 
		fullname=:f, 
		updated_at=now() 
		WHERE id=:i AND deleted_at IS NULL `, map[string]interface{}{
		"f": req.FullName,
		"i": req.Id,
	})
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n > 0 {
		return nil
	}
	return errors.New("author not found")
}

func (stg authorRepo) DeleteAuthor(id string) error {

	res, err := stg.db.Exec(`UPDATE author SET deleted_at=now() WHERE id=$1 AND deleted_at IS NULL`, id)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return errors.New("author not found")
	}
	return nil
}
