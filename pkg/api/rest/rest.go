package rest

import "github.com/iktzdx/skillfactory-gonews/pkg/api"

type PrimaryAdapter struct {
	port api.BoundaryPort
}

func NewPrimaryAdapter(port api.BoundaryPort) PrimaryAdapter {
	return PrimaryAdapter{port}
}
