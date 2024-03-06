package rest_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"

	"github.com/iktzdx/skillfactory-gonews/internal/app/posts"
	"github.com/iktzdx/skillfactory-gonews/pkg/api/rest"
	"github.com/iktzdx/skillfactory-gonews/pkg/storage"
)

type GetPostByIDSuite struct {
	suite.Suite
	req     *http.Request
	resp    *httptest.ResponseRecorder
	port    *rest.MockBoundaryPort
	adapter rest.PrimaryAdapter
}

func TestGetPostByIDSuite(t *testing.T) {
	suite.Run(t, new(GetPostByIDSuite))
}

func (s *GetPostByIDSuite) SetupTest() {
	req, err := http.NewRequest(http.MethodGet, "/post/12345", nil) //nolint:noctx
	s.Require().NoError(err, "make new get request")

	s.req = mux.SetURLVars(req, map[string]string{"id": "12345"})

	s.resp = httptest.NewRecorder()

	s.port = new(rest.MockBoundaryPort)
	s.adapter = rest.NewPrimaryAdapter(s.port)
}

func (s *GetPostByIDSuite) TestGetPostThatDoesNotExist() {
	var post posts.Post

	s.port.On("GetPostByID", "12345").Return(post, storage.ErrNoDataFound)
	s.adapter.GetPostByID(s.resp, s.req)
	s.Equal(http.StatusNotFound, s.resp.Code)

	var errMsg rest.WebAPIErrorResponse
	err := json.NewDecoder(s.resp.Body).Decode(&errMsg)
	s.Require().NoError(err, "decode web API error message")
	s.Equal(rest.PostNotFoundCode, errMsg.Code)
	s.Equal("no post found with id provided", errMsg.Message)
}

func (s *GetPostByIDSuite) TestGetPostThatDoesExist() {
	expected := posts.Post{
		ID:        12345,
		AuthorID:  0,
		Title:     "The Future of Sustainable Energy",
		Content:   "The global pursuit of renewable energy sources continues to gain momentum.",
		CreatedAt: 0,
	}

	s.port.On("GetPostByID", "12345").Return(expected, nil)
	s.adapter.GetPostByID(s.resp, s.req)
	s.Equal(http.StatusOK, s.resp.Code)

	want := `{
        "id": 12345,
        "authorId": 0,
        "title": "The Future of Sustainable Energy",
        "content": "The global pursuit of renewable energy sources continues to gain momentum.",
        "createdAt": 0
    }`

	s.JSONEq(want, s.resp.Body.String())
}

func (s *GetPostByIDSuite) TestGetPostWithInvalidID() {
	var post posts.Post

	s.port.On("GetPostByID", "12345").Return(post, posts.ErrInvalidQueryParam)
	s.adapter.GetPostByID(s.resp, s.req)
	s.Equal(http.StatusBadRequest, s.resp.Code)

	var errMsg rest.WebAPIErrorResponse
	err := json.NewDecoder(s.resp.Body).Decode(&errMsg)
	s.Require().NoError(err, "decode web API error message")
	s.Equal(rest.BadRequestCode, errMsg.Code)
	s.Equal("invalid query params provided", errMsg.Message)
}

func (s *GetPostByIDSuite) TestGetPostReturnsUnexpectedErr() {
	var post posts.Post

	s.port.On("GetPostByID", "12345").Return(post, posts.ErrUnexpected)
	s.adapter.GetPostByID(s.resp, s.req)
	s.Equal(http.StatusInternalServerError, s.resp.Code)

	var errMsg rest.WebAPIErrorResponse
	err := json.NewDecoder(s.resp.Body).Decode(&errMsg)
	s.Require().NoError(err, "decode web API error message")
	s.Equal(rest.UnexpectedCode, errMsg.Code)
	s.Equal("unexpected error attempting to get post", errMsg.Message)
}
