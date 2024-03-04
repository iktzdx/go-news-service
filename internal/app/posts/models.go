package posts

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
