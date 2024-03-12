package posts_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/iktzdx/skillfactory-gonews/internal/app/posts"
	"github.com/iktzdx/skillfactory-gonews/pkg/storage"
)

var errMockUnexpected = errors.New("unexpected")

type FindPostByIDSuite struct {
	suite.Suite
	mockRepo *posts.MockBoundaryRepoPort
	port     posts.BoundaryPort
}

func TestFindPostByIDSuite(t *testing.T) {
	suite.Run(t, new(FindPostByIDSuite))
}

func (s *FindPostByIDSuite) SetupTest() {
	s.mockRepo = new(posts.MockBoundaryRepoPort)
	s.port = posts.NewBoundaryPort(s.mockRepo)
}

func (s *FindPostByIDSuite) TestFinderFailed() {
	var expected storage.Data

	s.mockRepo.On("FindPostByID", int64(12345)).Return(expected, errMockUnexpected)

	got, err := s.port.GetPostByID("12345")

	s.Require().Error(err)
	s.Zero(got)
}

func (s *FindPostByIDSuite) TestInvalidPostID() {
	got, err := s.port.GetPostByID("12C45")

	s.Require().ErrorIs(err, posts.ErrInvalidQueryParam)
	s.Zero(got)
}

func (s *FindPostByIDSuite) TestFinderSucceeded() {
	expected := storage.Data{
		ID:        12345,
		AuthorID:  0,
		Title:     "The Future of Sustainable Energy",
		Content:   "The global pursuit of renewable energy sources continues to gain momentum.",
		CreatedAt: 0,
	}

	s.mockRepo.On("FindPostByID", expected.ID).Return(expected, nil)

	want := posts.FromRepo(expected)
	got, err := s.port.GetPostByID("12345")

	s.Require().NoError(err)
	s.Equal(want, got)
}
