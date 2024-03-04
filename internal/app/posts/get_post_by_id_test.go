package posts_test

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/iktzdx/skillfactory-gonews/internal/app/posts"
	"github.com/iktzdx/skillfactory-gonews/pkg/storage"
)

type FindPostByIDSuite struct {
	suite.Suite
	mockRepo *MockBoundaryRepoPort
	port     posts.BoundaryPort
}

func TestFindPostByIDSuite(t *testing.T) {
	suite.Run(t, new(FindPostByIDSuite))
}

type MockBoundaryRepoPort struct {
	mock.Mock
}

func (mockRepo *MockBoundaryRepoPort) FindPostByID(id int) (storage.Data, error) {
	args := mockRepo.Called(id)

	return args.Get(0).(storage.Data), args.Error(1) //nolint:forcetypeassert,wrapcheck
}

func (s *FindPostByIDSuite) SetupTest() {
	s.mockRepo = new(MockBoundaryRepoPort)
	s.port = posts.NewBoundaryPort(s.mockRepo)
}

func (s *FindPostByIDSuite) TestFinderFailed() {
	var expected storage.Data

	s.mockRepo.On("FindPostByID", 12345).Return(expected, storage.ErrUnexpected)

	want := posts.FromRepo(expected)
	got, err := s.port.GetPostByID("12345")

	s.Require().ErrorIs(err, storage.ErrUnexpected)
	s.Equal(want, got)
}

func (s *FindPostByIDSuite) TestInvalidPostID() {
	want := posts.Post{} //nolint:exhaustruct
	got, err := s.port.GetPostByID("12C45")

	s.Require().ErrorIs(err, posts.ErrInvalidQueryParam)
	s.Equal(want, got)
}

func (s *FindPostByIDSuite) TestFinderSucceeded() {
	expected := storage.Data{
		ID:        12345,
		AuthorID:  0,
		Title:     "The Future of Sustainable Energy",
		Content:   "The global pursuit of renewable energy sources continues to gain momentum.",
		CreatedAt: 0,
	}

	s.mockRepo.On("FindPostByID", 12345).Return(expected, nil)

	want := posts.FromRepo(expected)
	got, err := s.port.GetPostByID("12345")

	s.Require().NoError(err)
	s.Equal(want, got)
}
