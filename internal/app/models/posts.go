package models

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
