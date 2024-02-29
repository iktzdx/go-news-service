package rest

type PostsBoundaryPort interface {
	GetPost(id string) (Post, error)
	// List(opts SearchOpts) (Posts, error)
}
