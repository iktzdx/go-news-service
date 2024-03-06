package posts_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/iktzdx/skillfactory-gonews/internal/app/posts"
	"github.com/iktzdx/skillfactory-gonews/pkg/storage"
)

type ListPostsSuite struct {
	suite.Suite
	mockRepo *posts.MockBoundaryRepoPort
	port     posts.BoundaryPort
}

func TestListPostsSuite(t *testing.T) {
	suite.Run(t, new(ListPostsSuite))
}

func (s *ListPostsSuite) SetupTest() {
	s.mockRepo = new(posts.MockBoundaryRepoPort)
	s.port = posts.NewBoundaryPort(s.mockRepo)
}

func (s *ListPostsSuite) TestListPostsInvalidAuthorID() {
	var params posts.QueryParams

	params.AuthorID = "0xff"
	got, err := s.port.List(params)

	s.Require().ErrorIs(err, posts.ErrInvalidQueryParam)
	s.Zero(got)
}

func (s *ListPostsSuite) TestListPostsFailed() {
	var (
		expected storage.BulkData
		opts     storage.SearchOpts
		params   posts.QueryParams
	)

	s.mockRepo.On("List", opts).Return(expected, storage.ErrUnexpected)

	got, err := s.port.List(params)

	s.Require().ErrorIs(err, storage.ErrUnexpected)
	s.Zero(got)
}

func (s *ListPostsSuite) TestListPostsSucceeded() {
	expected := storage.BulkData{
		Posts: []storage.Data{
			{
				ID:        54321,
				AuthorID:  42,
				Title:     "This test is green #1",
				Content:   "The text data for green test",
				CreatedAt: 20120822193532,
			},
			{
				ID:        54322,
				AuthorID:  42,
				Title:     "This test is green #2",
				Content:   "The text data for green test",
				CreatedAt: 20120822193534,
			},
		},
		Total: 2,
	}

	var opts storage.SearchOpts
	opts.AuthorID = 42

	s.mockRepo.On("List", opts).Return(expected, nil)

	var params posts.QueryParams
	params.AuthorID = "42"

	want := posts.FromRepoBulk(expected)
	got, err := s.port.List(params)

	s.Require().NoError(err)
	s.Equal(want, got)
}
