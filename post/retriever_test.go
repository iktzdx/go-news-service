package post_test

import (
	"errors"
	"gonews/api"
	"gonews/post"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

var ErrUnexpected = errors.New("unexpected error")

type PostRetrieverSuite struct {
	suite.Suite
	m    *MockPostRetriever
	port post.PostRetriever
}

func TestPostRetrieverSuite(t *testing.T) {
	suite.Run(t, new(PostRetrieverSuite))
}

type MockPostRetriever struct {
	mock.Mock
}

func (m *MockPostRetriever) FindPostByID(id int) (api.Post, error) {
	args := m.Called(id)

	return args.Get(0).(api.Post), args.Error(1)
}

func (s *PostRetrieverSuite) SetupTest() {
	s.m = new(MockPostRetriever)
	s.port = post.NewPostRetriever(s.m)
}

func (s *PostRetrieverSuite) TestRetrieverFailed() {
	var p api.Post

	s.m.On("FindPostByID", 12345).Return(p, ErrUnexpected)

	_, err := s.port.GetPost(12345)
	s.Require().ErrorIs(err, ErrUnexpected)
}

func (s *PostRetrieverSuite) TestRetrieverSuccess() {
	expected := api.Post{
		ID:        12345,
		AuthorID:  0,
		Title:     "The Future of Sustainable Energy",
		Content:   "The global pursuit of renewable energy sources continues to gain momentum.",
		CreatedAt: 0,
	}

	s.m.On("FindPostByID", 12345).Return(expected, nil)

	got, err := s.port.GetPost(12345)
	s.Require().NoError(err)
	s.Equal(expected, got)
}
