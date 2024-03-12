package storage

type BulkData struct {
	Posts []Data
	Total int64
}

type Data struct {
	ID        int64
	AuthorID  int64
	Title     string
	Content   string
	CreatedAt int64
}

type (
	SearchOpts struct {
		FilterOpts
		PaginationOpts
	}

	FilterOpts struct {
		ID       int64
		AuthorID int64
	}

	PaginationOpts struct {
		Limit  int
		Offset int
	}
)
