package pgsql

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/iktzdx/skillfactory-gonews/pkg/storage"
)

const defaultLimit int = 5

type SecondaryAdapter struct {
	db *sql.DB
}

func NewSecondaryAdapter(db *sql.DB) SecondaryAdapter {
	return SecondaryAdapter{db}
}

func (adapter SecondaryAdapter) Create(data storage.Data) (int64, error) {
	query := `INSERT INTO posts (author_id, title, content, created_at)
    VALUES($1, $2, $3, $4) RETURNING id`

	var createdID int64

	row := adapter.db.QueryRow(query, data.AuthorID, data.Title, data.Content, data.CreatedAt) //nolint:execinquery
	if err := row.Scan(&createdID); err != nil {
		return -1, errors.Wrap(err, "scan for the last inserted ID")
	}

	return createdID, nil
}

func (adapter SecondaryAdapter) FindPostByID(id int64) (storage.Data, error) {
	var post storage.Data

	query := "SELECT * FROM posts WHERE id = $1"

	row := adapter.db.QueryRow(query, id)
	if err := row.Scan(&post.ID, &post.AuthorID, &post.Title, &post.Content, &post.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return storage.Data{}, storage.ErrNoDataFound
		}

		return storage.Data{}, errors.Wrap(err, "scan query row")
	}

	return post, nil
}

func (adapter SecondaryAdapter) List(opts storage.SearchOpts) (storage.BulkData, error) {
	query := `SELECT * FROM posts WHERE ($1 = 0 OR id = $1) AND ($2 = 0 OR author_id = $2) LIMIT $3 OFFSET $4`

	limit := opts.Limit
	if limit == 0 {
		limit = defaultLimit
	}

	rows, err := adapter.db.Query(query, opts.ID, opts.AuthorID, limit, opts.Offset)
	if err != nil {
		return storage.BulkData{}, errors.Wrap(err, "query rows")
	}

	defer rows.Close()

	if err := rows.Err(); err != nil {
		return storage.BulkData{}, errors.Wrap(err, "check rows err")
	}

	var bulkData storage.BulkData

	for rows.Next() {
		var data storage.Data
		if err := rows.Scan(&data.ID, &data.AuthorID, &data.Title, &data.Content, &data.CreatedAt); err != nil {
			return storage.BulkData{}, errors.Wrap(err, "scan data")
		}

		bulkData.Posts = append(bulkData.Posts, data)
		bulkData.Total++
	}

	return bulkData, nil
}

func (adapter SecondaryAdapter) Update(change storage.Data) (int64, error) {
	query := `
    UPDATE posts
    SET author_id = $2, title = $3, content = $4
    WHERE id = $1`

	row, err := adapter.db.Exec(query, change.ID, change.AuthorID, change.Title, change.Content)
	if err != nil {
		return -1, errors.Wrap(err, "update row")
	}

	affected, err := row.RowsAffected()
	if err != nil {
		return -1, errors.Wrap(err, "get rows affected amount")
	}

	if affected == 0 {
		return -1, storage.ErrNoDataFound
	}

	return affected, nil
}

func (adapter SecondaryAdapter) Delete(id int64) (int64, error) {
	query := `DELETE FROM posts WHERE id = $1`

	row, err := adapter.db.Exec(query, id)
	if err != nil {
		return -1, errors.Wrap(err, "delete row")
	}

	affected, err := row.RowsAffected()
	if err != nil {
		return -1, errors.Wrap(err, "get rows affected amount")
	}

	if affected == 0 {
		return -1, storage.ErrNoDataFound
	}

	return affected, nil
}
