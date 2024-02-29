package rest_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/iktzdx/skillfactory-gonews/internal/app/rest"
)

type GetPostByIDSuite struct {
	suite.Suite
	req     *http.Request
	resp    *httptest.ResponseRecorder
	port    *MockBoundaryPort
	adapter rest.PrimaryAdapter
}

func TestGetPostByIDSuite(t *testing.T) {
	suite.Run(t, new(GetPostByIDSuite))
}

type MockBoundaryPort struct {
	mock.Mock
}

func (m *MockBoundaryPort) GetPostByID(id string) (rest.Post, error) {
	args := m.Called(id)

	return args.Get(0).(rest.Post), args.Error(1) //nolint:forcetypeassert,wrapcheck
}

func (s *GetPostByIDSuite) SetupTest() {
	req, err := http.NewRequest(http.MethodGet, "/post/12345", nil) //nolint:noctx
	s.Require().NoError(err, "make new get request")

	s.req = mux.SetURLVars(req, map[string]string{"id": "12345"})

	s.resp = httptest.NewRecorder()

	s.port = new(MockBoundaryPort)
	s.adapter = rest.NewPrimaryAdapter(s.port)
}

func (s *GetPostByIDSuite) TestGetPostThatDoesNotExist() {
	var post rest.Post

	s.port.On("GetPostByID", "12345").Return(post, rest.ErrPostNotFound)
	s.adapter.GetPostByID(s.resp, s.req)
	s.Equal(http.StatusNotFound, s.resp.Code)

	var errMsg rest.WebAPIError
	err := json.NewDecoder(s.resp.Body).Decode(&errMsg)
	s.Require().NoError(err, "decode web API error message")
	s.Equal("001", errMsg.Code)
	s.Equal("no post with id 12345", errMsg.Message)
}

func (s *GetPostByIDSuite) TestGetPostThatDoesExist() {
	post := rest.Post{
		ID:        12345,
		AuthorID:  0,
		Title:     "The Future of Sustainable Energy",
		Content:   "The global pursuit of renewable energy sources continues to gain momentum.",
		CreatedAt: 0,
	}

	s.port.On("GetPostByID", "12345").Return(post, nil)
	s.adapter.GetPostByID(s.resp, s.req)
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

func (s *GetPostByIDSuite) TestGetPostWithInvalidID() {
	var post rest.Post

	s.port.On("GetPostByID", "12345").Return(post, rest.ErrInvalidPostID)
	s.adapter.GetPostByID(s.resp, s.req)
	s.Equal(http.StatusBadRequest, s.resp.Code)

	var errMsg rest.WebAPIError
	err := json.NewDecoder(s.resp.Body).Decode(&errMsg)
	s.Require().NoError(err, "decode web API error message")
	s.Equal("003", errMsg.Code)
	s.Equal("invalid post id provided", errMsg.Message)
}

func (s *GetPostByIDSuite) TestGetPostReturnsUnexpectedErr() {
	var post rest.Post

	s.port.On("GetPostByID", "12345").Return(post, rest.ErrUnexpected)
	s.adapter.GetPostByID(s.resp, s.req)
	s.Equal(http.StatusInternalServerError, s.resp.Code)

	var errMsg rest.WebAPIError
	err := json.NewDecoder(s.resp.Body).Decode(&errMsg)
	s.Require().NoError(err, "decode web API error message")
	s.Equal("002", errMsg.Code)
	s.Equal("unexpected error attempting to get post", errMsg.Message)
}
