package api_test

import (
	"context"
	"encoding/json"
	"gonews/api"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type GetPostSuite struct {
	suite.Suite
	req  *http.Request
	resp *httptest.ResponseRecorder
	r    *MockPostRetriever
	h    api.GetPostHandler
}

func TestGetPostSuite(t *testing.T) {
	suite.Run(t, new(GetPostSuite))
}

type MockPostRetriever struct {
	mock.Mock
}

func (m *MockPostRetriever) GetPost(id int) (api.Post, error) {
	args := m.Called(id)

	return args.Get(0).(api.Post), args.Error(1)
}

func (s *GetPostSuite) SetupTest() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/post/12345", nil)
	s.Require().NoError(err, "make new get request")

	s.req = mux.SetURLVars(req, map[string]string{"id": "12345"})

	s.resp = httptest.NewRecorder()

	s.r = new(MockPostRetriever)
	s.h = api.NewGetPostHandler(s.r)
}

func (s *GetPostSuite) TestGetPostThatDoesNotExist() {
	var post api.Post

	s.r.On("GetPost", 12345).Return(post, api.ErrPostNotFound)
	s.h.ServeHTTP(s.resp, s.req)
	s.Equal(http.StatusNotFound, s.resp.Code)

	var errMsg api.WebAPIError
	err := json.NewDecoder(s.resp.Body).Decode(&errMsg)
	s.Require().NoError(err, "decode web API error message")
	s.Equal("001", errMsg.Code)
	s.Equal("no post with id 12345", errMsg.Message)
}

func (s *GetPostSuite) TestGetPostThatDoesExist() {
	post := api.Post{
		ID:        12345,
		AuthorID:  0,
		Title:     "The Future of Sustainable Energy",
		Content:   "The global pursuit of renewable energy sources continues to gain momentum.",
		CreatedAt: 0,
	}

	s.r.On("GetPost", 12345).Return(post, nil)

	s.h.ServeHTTP(s.resp, s.req)
	s.Equal(http.StatusOK, s.resp.Code)

	expectedBody := `{
        "id": 12345,
        "authorId": 0,
        "title": "The Future of Sustainable Energy",
        "content": "The global pursuit of renewable energy sources continues to gain momentum.",
        "createdAt": 0
    }`

	s.JSONEq(expectedBody, s.resp.Body.String())
}

func (s *GetPostSuite) TestGetPostReturnsUnexpectedErr() {
	var post api.Post

	s.r.On("GetPost", 12345).Return(post, api.ErrUnexpected)

	s.h.ServeHTTP(s.resp, s.req)

	s.Equal(http.StatusInternalServerError, s.resp.Code)

	var errMsg api.WebAPIError
	err := json.NewDecoder(s.resp.Body).Decode(&errMsg)
	s.Require().NoError(err, "decode web API error message")
	s.Equal("002", errMsg.Code)
	s.Equal("unexpected error attempting to get post", errMsg.Message)
}
