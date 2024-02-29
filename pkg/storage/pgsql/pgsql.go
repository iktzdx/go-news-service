package pgsql

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/iktzdx/skillfactory-gonews/internal/app/models"
	"github.com/iktzdx/skillfactory-gonews/pkg/api/rest"
)

type SecondaryAdapter struct {
	db *sql.DB
}

func NewSecondaryAdapter(db *sql.DB) SecondaryAdapter {
	return SecondaryAdapter{db}
}

func (adapter SecondaryAdapter) FindPostByID(id int) (models.Post, error) {
	var post models.Post

	query := "SELECT * FROM posts WHERE id = $1"
	row := adapter.db.QueryRow(query, id)

	if err := row.Scan(&post.ID, &post.AuthorID, &post.Title, &post.Content, &post.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Post{}, rest.ErrPostNotFound
		}

		return models.Post{}, rest.ErrUnexpected
	}

	return post, nil
}
