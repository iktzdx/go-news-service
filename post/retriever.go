package post

import (
	"database/sql"
	"errors"
	"gonews/api"

	_ "github.com/lib/pq"
)

type PostRetriever struct {
	db *sql.DB
}

func NewPostRetriever(db *sql.DB) PostRetriever {
	return PostRetriever{db}
}

func (r PostRetriever) GetPost(id int) (api.Post, error) {
	var post api.Post

	query := "SELECT * FROM posts WHERE id = $1"
	row := r.db.QueryRow(query, id)

	if err := row.Scan(&post.ID, &post.AuthorID, &post.Title, &post.Content, &post.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return api.Post{}, api.ErrPostNotFound
		}

		return api.Post{}, api.ErrUnexpected
	}

	return post, nil
}
