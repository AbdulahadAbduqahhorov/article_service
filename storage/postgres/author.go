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
func (u authorRepo) CreateAuthor(req *author_service.CreateAuthorRequest) (string, error) {
	id := uuid.New().String()
	_, err := u.db.Exec(`INSERT INTO 
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

func (u authorRepo) GetAuthor(req *author_service.GetAuthorRequest) (*author_service.GetAuthorResponse, error) {
	res := &author_service.GetAuthorResponse{
		Authors: make([]*author_service.Author, 0),
	}
	var (
		updatedAt sql.NullString
		deletedAt sql.NullString
	)
	rows, err := u.db.Queryx(`SELECT 
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
			&author.FullName,
			&author.CreatedAt,
			&updatedAt,
			&deletedAt,
		)
		if err != nil {
			return nil, err
		}
		if updatedAt.Valid {
			author.UpdatedAt = updatedAt.String
		}
		if deletedAt.Valid {
			author.DeletedAt = deletedAt.String
		}
		res.Authors = append(res.Authors, &author)

	}

	return res, nil

}

func (u authorRepo) GetAuthorById(id string) (*author_service.Author, error) {
	res:=&author_service.Author{}
	var (
		updatedAt sql.NullString
		deletedAt sql.NullString
	)
	err := u.db.QueryRow(`
	SELECT 
		id,
		fullname,
		created_at,
		updated_at,
		deleted_at
	FROM author  
	WHERE id=$1 AND deleted_at is NULL`, id).Scan(
		&res.Id,
		&res.FullName,
		&res.CreatedAt,
		&updatedAt,
		&deletedAt,
	)
	if updatedAt.Valid {
		res.UpdatedAt = updatedAt.String
	}
	if deletedAt.Valid {
		res.DeletedAt = deletedAt.String
	}
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (u authorRepo) UpdateAuthor(req *author_service.UpdateAuthorRequest) error {
	res, err := u.db.NamedExec(`
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

func (u authorRepo) DeleteAuthor(id string) error {

	res, err := u.db.Exec(`UPDATE author SET deleted_at=now() WHERE id=$1 AND deleted_at IS NULL`, id)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n > 0 {
		return nil
	}
	return errors.New("author not found")
}
