package postgres

import (
	"database/sql"
	"errors"
	"time"

	"github.com/AbdulahadAbduqahhorov/gin/Article/genproto/author_service"
	"github.com/AbdulahadAbduqahhorov/gin/Article/models"
	"github.com/AbdulahadAbduqahhorov/gin/Article/storage"
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
func (stg authorRepo) CreateAuthor(req author_service.CreateAuthorRequest) (author_service.CreateAuthorResponse,error) {
	id:=uuid.New().String()
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
		return author_service.CreateAuthorResponse{},err
	}
	return author_service.CreateAuthorResponse{Id:id},nil
	
}

func (stg authorRepo) GetAuthor(req author_service.GetAuthorRequest) (author_service.GetAuthorResponse, error) {
	var res author_service.GetAuthorResponse
	var tempFullname *string

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
		req.Limit,
		req.Search,
		req.Offset,
	)
	if err != nil {
		return res, err
	}

	for rows.Next() {
		var author author_service.Author
		err := rows.Scan(
			&author.Id,
			&tempFullname,
			&author.CreatedAt,
			&author.UpdatedAt,
			&author.DeletedAt,
		)
		if err != nil {
			return res, err
		}
		if tempFullname != nil {
			author.FullName = *tempFullname
		}
		res.Authors = append(res.Authors, &author)

	}

	return res, err

}

func (stg authorRepo) GetAuthorById(req author_service.GetAuthorByIdResponse) (author_service.Author, error) {
	var res  author_service.Author
	var tempFullname *string
	var tempUpdatedAt *time.Time

	var a sql.NullString
	err := stg.db.QueryRow(`
	SELECT 
		id,
		fullname,
		created_at,
		updated_at,
		deleted_at
	FROM author  
	WHERE id=$1 AND deleted_at is NULL`, req.Id).Scan(
		&res.Id,
		&tempFullname,
		&res.CreatedAt,
		&a,
		&res.DeletedAt,
	)
	if a.Valid{
		res.UpdatedAt=a.String
	}
	if err != nil {
		return res, err
	}
	if tempUpdatedAt!=nil{
		res.UpdatedAt=
	}
	if tempFullname != nil {
		res.FullName = *tempFullname
	}
	return res, nil
}

func (stg authorRepo) UpdateAuthor(author models.UpdateAuthorModel) error {
	res, err := stg.db.NamedExec(`
	UPDATE  author SET 
		fullname=:f, 
		updated_at=now() 
		WHERE id=:i AND deleted_at IS NULL `, map[string]interface{}{
		"f": author.FullName,
		"i": author.Id,
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
