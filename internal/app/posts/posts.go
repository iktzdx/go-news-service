package posts

import "github.com/iktzdx/skillfactory-gonews/pkg/storage"

type BoundaryPort struct {
	repo storage.BoundaryRepoPort
}

func NewBoundaryPort(repo storage.BoundaryRepoPort) BoundaryPort {
	return BoundaryPort{repo}
}
