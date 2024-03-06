package rest_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/iktzdx/skillfactory-gonews/internal/app/posts"
	"github.com/iktzdx/skillfactory-gonews/pkg/api/rest"
)

type ListPostsSuite struct {
	suite.Suite
	req     *http.Request
	resp    *httptest.ResponseRecorder
	port    *rest.MockBoundaryPort
	adapter rest.PrimaryAdapter
}

func TestListPostsSuite(t *testing.T) {
	suite.Run(t, new(ListPostsSuite))
}

func (s *ListPostsSuite) SetupTest() {
	req, err := http.NewRequest(http.MethodGet, "/posts", nil)
	s.Require().NoError(err, "make new get request")

	q := req.URL.Query()
	req.URL.RawQuery = q.Encode()

	s.req = req
	s.resp = httptest.NewRecorder()

	s.port = new(rest.MockBoundaryPort)
	s.adapter = rest.NewPrimaryAdapter(s.port)
}

func (s *ListPostsSuite) TestListPostsEmptyResult() {
	params := posts.QueryParams{
		FiltersParams: posts.FiltersParams{
			ID:       "12345",
			AuthorID: "12345",
		},
		PaginationParams: posts.PaginationParams{
			Limit:  "10",
			Offset: "0",
		},
	}

	s.port.On("List", params).Return(posts.Posts{}, nil) //nolint:exhaustruct

	q := s.req.URL.Query()
	q.Add("id", "12345")
	q.Add("author_id", "12345")
	q.Add("limit", "10")
	q.Add("offset", "0")
	s.req.URL.RawQuery = q.Encode()

	s.adapter.List(s.resp, s.req)
	s.Equal(http.StatusOK, s.resp.Code)

	want := `{"posts": [], "total": 0}`
	s.JSONEq(want, s.resp.Body.String())
}

func (s *ListPostsSuite) TestListPostsInvalidAuthorID() {
	var params posts.QueryParams
	params.AuthorID = "0xFF"

	s.port.On("List", params).Return(posts.Posts{}, posts.ErrInvalidQueryParam) //nolint:exhaustruct

	q := s.req.URL.Query()
	q.Add("author_id", "0xFF")
	s.req.URL.RawQuery = q.Encode()

	s.adapter.List(s.resp, s.req)
	s.Equal(http.StatusBadRequest, s.resp.Code)

	var errMsg rest.WebAPIErrorResponse
	err := json.NewDecoder(s.resp.Body).Decode(&errMsg)
	s.Require().NoError(err, "decode web API error message")
	s.Equal(rest.BadRequestCode, errMsg.Code)
	s.Equal("invalid query params provided", errMsg.Message)
}

func (s *ListPostsSuite) TestListPostsReturnsUnexpectedErr() {
	s.port.On("List", posts.QueryParams{}).Return(posts.Posts{}, posts.ErrUnexpected) //nolint:exhaustruct

	s.adapter.List(s.resp, s.req)
	s.Equal(http.StatusInternalServerError, s.resp.Code)

	var errMsg rest.WebAPIErrorResponse
	err := json.NewDecoder(s.resp.Body).Decode(&errMsg)
	s.Require().NoError(err, "decode web API error message")
	s.Equal(rest.UnexpectedCode, errMsg.Code)
	s.Equal("service returned unexpected error", errMsg.Message)
}

func (s *ListPostsSuite) TestListPostsTheDoExist() {
	//nolint:exhaustruct
	params := posts.QueryParams{
		FiltersParams: posts.FiltersParams{
			AuthorID: "123",
		},
		PaginationParams: posts.PaginationParams{
			Offset: "1",
		},
	}

	expected := posts.Posts{
		Posts: []posts.Post{
			{
				ID:        54321,
				AuthorID:  123,
				Title:     "This test is green",
				Content:   "The content for the green test",
				CreatedAt: 20240307193532,
			},
		},
		Total: 1,
	}

	s.port.On("List", params).Return(expected, nil)

	q := s.req.URL.Query()
	q.Add("author_id", "123")
	q.Add("offset", "1")
	s.req.URL.RawQuery = q.Encode()

	s.adapter.List(s.resp, s.req)
	s.Equal(http.StatusOK, s.resp.Code)

	want := `{
        "posts":
            [{
                "id": 54321,
                "authorId": 123,
                "title": "This test is green",
                "content": "The content for the green test",
                "createdAt": 20240307193532
            }],
        "total": 1
    }`

	s.JSONEq(want, s.resp.Body.String())
}
