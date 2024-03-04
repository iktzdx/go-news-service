package storage

type BulkData struct {
	Posts []Data
}

type Data struct {
	ID        int
	AuthorID  int
	Title     string
	Content   string
	CreatedAt int
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
