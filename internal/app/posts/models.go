package posts

import "github.com/iktzdx/skillfactory-gonews/pkg/storage"

type Posts struct {
	Posts []Post
	Total int
}

type Post struct {
	ID        int64
	AuthorID  int64
	Title     string
	Content   string
	CreatedAt int64
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

func FromRepoBulk(bulkData storage.BulkData) Posts {
	total := len(bulkData.Posts)
	posts := make([]Post, total)

	for idx, data := range bulkData.Posts {
		posts[idx] = FromRepo(data)
	}

	return Posts{
		Posts: posts,
		Total: total,
	}
}
