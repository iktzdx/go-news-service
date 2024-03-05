package posts

import "github.com/iktzdx/skillfactory-gonews/pkg/storage"

type Posts struct {
	Posts []Post
	Total int
}

type Post struct {
	ID        int
	AuthorID  int
	Title     string
	Content   string
	CreatedAt int
}

type (
	QueryParams struct {
		FiltersParams
		PaginationParams
	}

	FiltersParams struct {
		ID       string
		AuthorID string
	}

	PaginationParams struct {
		Limit  string
		Offset string
	}
)

func FromRepo(data storage.Data) Post {
	return Post{
		ID:        data.ID,
		AuthorID:  data.AuthorID,
		Title:     data.Title,
		Content:   data.Content,
		CreatedAt: data.CreatedAt,
	}
}
