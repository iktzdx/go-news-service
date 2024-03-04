package rest

type (
	GetPostByIDResponse struct {
		Payload
	}

	ListResponse struct {
		Posts []Payload `json:"posts"`
		Total int       `json:"total"`
	}

	WebAPIErrorResponse struct {
		Code    string `json:"code"`
		Message string `json:"msg"`
	}
)

type Payload struct {
	ID        int    `json:"id"`
	AuthorID  int    `json:"authorId"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	CreatedAt int    `json:"createdAt"`
}
