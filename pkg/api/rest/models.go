package rest

type (
	GetPostByIDResponse struct {
		Payload
	}

	ListPostsResponse struct {
		Posts []Payload `json:"posts"`
		Total int       `json:"total"`
	}

	CreatePostResponse struct {
		ID int64 `json:"id"`
	}

	WebAPIErrorResponse struct {
		Code    string `json:"code"`
		Message string `json:"msg"`
	}
)

type (
	CreatePostRequest struct {
		AuthorID  int64  `json:"authorId"`
		Title     string `json:"title"`
		Content   string `json:"content"`
		CreatedAt int64  `json:"createdAt"`
	}
)

type Payload struct {
	ID        int64  `json:"id"`
	AuthorID  int64  `json:"authorId"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	CreatedAt int64  `json:"createdAt"`
}
