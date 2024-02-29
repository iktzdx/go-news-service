package storage

import "github.com/iktzdx/skillfactory-gonews/internal/app/rest"

type BoundaryRepoPort interface {
	// Create(post Post) (Post, error)
	FindPostByID(id int) (rest.Post, error)
	// List(opts SearchOpts) (Posts, error)
	// Update(post Post) (Post, error)
	// Delete(id string) error
}
