package api_test

import (
	"gonews/api"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type GetPostSuite struct {
	suite.Suite
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

func (s *GetPostSuite) TestGetPostThatDoesNotExist() {
	req, err := http.NewRequest(http.MethodGet, "/post/12345", nil)
	s.Require().NoError(err, "make new get request")

	req = mux.SetURLVars(req, map[string]string{"id": "12345"})

	resp := httptest.NewRecorder()

	var post api.Post

	r := new(MockPostRetriever)
	r.On("GetPost", 12345).Return(post, api.ErrPostNotFound)

	h := api.NewGetPostHandler(r)
	h.ServeHTTP(resp, req)

	body, err := io.ReadAll(resp.Body)
	s.Require().NoError(err, "read response body")

	s.Equal(http.StatusNotFound, resp.Code)

	expectedBody := `{"code": "001", "msg": "no post with id 12345"}`
	s.JSONEq(expectedBody, string(body))
}

func (s *GetPostSuite) TestGetPostThatDoesExist() {
	req, err := http.NewRequest(http.MethodGet, "/post/12345", nil)
	s.Require().NoError(err, "make new get request")

	req = mux.SetURLVars(req, map[string]string{"id": "12345"})

	resp := httptest.NewRecorder()

	post := api.Post{
		ID:        12345,
		AuthorID:  0,
		Title:     "The Future of Sustainable Energy",
		Content:   "The global pursuit of renewable energy sources continues to gain momentum.",
		CreatedAt: 0,
	}

	r := new(MockPostRetriever)
	r.On("GetPost", 12345).Return(post, nil)

	h := api.NewGetPostHandler(r)
	h.ServeHTTP(resp, req)

	body, err := io.ReadAll(resp.Body)
	s.Require().NoError(err, "read response body")

	expectedBody := `{
        "id": 12345,
        "authorId": 0,
        "title": "The Future of Sustainable Energy",
        "content": "The global pursuit of renewable energy sources continues to gain momentum.",
        "createdAt": 0
    }`

	s.Equal(http.StatusOK, resp.Code)
	s.JSONEq(expectedBody, string(body))
}

func (s *GetPostSuite) TestGetPostReturnsUnexpectedErr() {
	req, err := http.NewRequest(http.MethodGet, "/post/12345", nil)
	s.Require().NoError(err, "make new get request")

	req = mux.SetURLVars(req, map[string]string{"id": "12345"})

	resp := httptest.NewRecorder()

	var post api.Post

	r := new(MockPostRetriever)
	r.On("GetPost", 12345).Return(post, api.ErrUnexpected)

	h := api.NewGetPostHandler(r)
	h.ServeHTTP(resp, req)

	body, err := io.ReadAll(resp.Body)
	s.Require().NoError(err, "read response body")

	s.Equal(http.StatusInternalServerError, resp.Code)

	expectedBody := `{"code": "002", "msg": "unexpected error attempting to get post"}`
	s.JSONEq(expectedBody, string(body))
}
