package rest

type Posts struct {
	Posts []Post `json:"posts"`
	Total int    `json:"total"`
}

type Post struct {
	ID        int    `json:"id"`
	AuthorID  int    `json:"authorId"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	CreatedAt int    `json:"createdAt"`
}

type WebAPIError struct {
	Code    string `json:"code"`
	Message string `json:"msg"`
}

type (
	SearchOpts struct {
		FilterOpts
		PaginationOpts
	}
	FilterOpts struct {
		ID       int
		AuthorID int
	}
	PaginationOpts struct {
		Limit  int
		Offset int
	}
)
