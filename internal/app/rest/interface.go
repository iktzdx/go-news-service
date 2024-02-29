package rest

type BoundaryPort interface {
	// Create(post Post) (Post, error)
	GetPostByID(id string) (Post, error)
	// List(opts SearchOpts) (Posts, error)
	// Update(post Post) (Post, error)
	// Delete(id string) error
}
