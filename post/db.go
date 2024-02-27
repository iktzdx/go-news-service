package post

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"

	"host.local/gonews/api"
)

type PGSQLSecondaryAdapter struct {
	db *sql.DB
}

func NewPGSQLSecondaryAdapter(db *sql.DB) PGSQLSecondaryAdapter {
	return PGSQLSecondaryAdapter{db}
}

func (adapter PGSQLSecondaryAdapter) FindPostByID(id int) (api.Post, error) {
	var post api.Post

	query := "SELECT * FROM posts WHERE id = $1"
	row := adapter.db.QueryRow(query, id)

	if err := row.Scan(&post.ID, &post.AuthorID, &post.Title, &post.Content, &post.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return api.Post{}, api.ErrPostNotFound
		}

		return api.Post{}, api.ErrUnexpected
	}

	return post, nil
}
