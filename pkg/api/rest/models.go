package rest

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
