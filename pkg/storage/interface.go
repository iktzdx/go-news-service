package storage

type BoundaryRepoPort interface {
	// Create(post Post) (Post, error)
	FindPostByID(id int) (Data, error)
	List(opts SearchOpts) (BulkData, error)
	// Update(post Post) (Post, error)
	// Delete(id string) error
}
