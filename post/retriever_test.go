package post_test

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"host.local/gonews/api"
	"host.local/gonews/post"
)

type FindPostByIDSuite struct {
	suite.Suite
	mockRepo *MockPostsBoundaryRepoPort
	port     post.PostsBoundaryPort
}

func TestFindPostByIDSuite(t *testing.T) {
	suite.Run(t, new(FindPostByIDSuite))
}

type MockPostsBoundaryRepoPort struct {
	mock.Mock
}

func (mockRepo *MockPostsBoundaryRepoPort) FindPostByID(id int) (api.Post, error) {
	args := mockRepo.Called(id)

	return args.Get(0).(api.Post), args.Error(1) //nolint:forcetypeassert,wrapcheck
}

func (s *FindPostByIDSuite) SetupTest() {
	s.mockRepo = new(MockPostsBoundaryRepoPort)
	s.port = post.NewPostsBoundaryPort(s.mockRepo)
}

func (s *FindPostByIDSuite) TestFinderFailed() {
	var expected api.Post

	s.mockRepo.On("FindPostByID", 12345).Return(expected, api.ErrUnexpected)

	got, err := s.port.GetPost("12345")
	s.Require().ErrorIs(err, api.ErrUnexpected)
	s.Equal(expected, got)
}

func (s *FindPostByIDSuite) TestInvalidPostID() {
	var expected api.Post

	s.mockRepo.On("FindPostByID", 12345).Return(expected, api.ErrInvalidPostID)

	got, err := s.port.GetPost("12C45")
	s.Require().ErrorIs(err, api.ErrInvalidPostID)
	s.Equal(expected, got)
}

func (s *FindPostByIDSuite) TestFinderSucceeded() {
	expected := api.Post{
		ID:        12345,
		AuthorID:  0,
		Title:     "The Future of Sustainable Energy",
		Content:   "The global pursuit of renewable energy sources continues to gain momentum.",
		CreatedAt: 0,
	}

	s.mockRepo.On("FindPostByID", 12345).Return(expected, nil)

	got, err := s.port.GetPost("12345")
	s.Require().NoError(err)
	s.Equal(expected, got)
}
