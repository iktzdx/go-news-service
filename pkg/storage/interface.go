package storage

type BoundaryRepoPort interface {
	Create(data Data) (int64, error)
	FindPostByID(id int64) (Data, error)
	List(opts SearchOpts) (BulkData, error)
	Update(change Data) (int64, error)
	Delete(id int64) (int64, error)
}
