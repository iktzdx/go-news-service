package rest

import (
	"github.com/stretchr/testify/mock"

	"github.com/iktzdx/skillfactory-gonews/internal/app/posts"
)

type MockBoundaryPort struct {
	mock.Mock
}

func (m *MockBoundaryPort) GetPostByID(id string) (posts.Post, error) {
	args := m.Called(id)

	return args.Get(0).(posts.Post), args.Error(1) //nolint:forcetypeassert,wrapcheck
}

func (m *MockBoundaryPort) List(params posts.QueryParams) (posts.Posts, error) {
	args := m.Called(params)

	return args.Get(0).(posts.Posts), args.Error(1) //nolint:forcetypeassert,wrapcheck
}
