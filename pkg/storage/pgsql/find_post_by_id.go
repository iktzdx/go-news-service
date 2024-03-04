package pgsql

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/iktzdx/skillfactory-gonews/pkg/storage"
)

type SecondaryAdapter struct {
	db *sql.DB
}

func NewSecondaryAdapter(db *sql.DB) SecondaryAdapter {
	return SecondaryAdapter{db}
}

func (adapter SecondaryAdapter) FindPostByID(id int) (storage.Data, error) {
	var post storage.Data

	query := "SELECT * FROM posts WHERE id = $1"
	row := adapter.db.QueryRow(query, id)

	if err := row.Scan(&post.ID, &post.AuthorID, &post.Title, &post.Content, &post.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return storage.Data{}, storage.ErrNoDataFound
		}

		return storage.Data{}, storage.ErrUnexpected
	}

	return post, nil
}
