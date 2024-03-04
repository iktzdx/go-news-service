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
