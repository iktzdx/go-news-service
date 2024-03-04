package test_test

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/iktzdx/skillfactory-gonews/pkg/api/rest"
)

type ListPostsSuite struct {
	suite.Suite
	c *http.Client
}

func TestListPostsSuite(t *testing.T) {
	suite.Run(t, new(ListPostsSuite))
}

func (s *ListPostsSuite) SetupTest() {
	s.c = http.DefaultClient
}

func (s *ListPostsSuite) TestListPostsFilteredNoResult() {
	req, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:8080/posts", nil)
	s.Require().NoError(err, "create new get request")

	q := req.URL.Query()
	q.Add("id", "12345")
	q.Add("author_id", "42")

	req.URL.RawQuery = q.Encode()

	resp, err := s.c.Do(req)
	s.Require().NoError(err, "make the actual get request")

	defer resp.Body.Close()

	s.Equal(http.StatusOK, resp.StatusCode)

	b, err := io.ReadAll(resp.Body)
	s.Require().NoError(err, "read response body")

	want := `{"posts": [], "total": 0}`
	s.JSONEq(want, string(b))
}

func (s *ListPostsSuite) TestListPostsFiltered() {
	req, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:8080/posts", nil) //nolint:noctx
	s.Require().NoError(err, "create new get request")

	q := req.URL.Query()
	q.Add("id", "0")
	q.Add("author_id", "0")

	req.URL.RawQuery = q.Encode()

	resp, err := s.c.Do(req)
	s.Require().NoError(err, "make the actual get request")

	defer resp.Body.Close()

	s.Equal(http.StatusOK, resp.StatusCode)

	b, err := io.ReadAll(resp.Body)
	s.Require().NoError(err, "read response body")

	want := `{"posts": [{"id": 0, "authorId": 0, "title": "Default Post", "content": "Simple text data", "createdAt": 0}], "total": 1}`
	s.JSONEq(want, string(b))
}

func (s *ListPostsSuite) TestListPostsInvalidAuthorID() {
	req, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:8080/posts", nil) //nolint:noctx
	s.Require().NoError(err, "create new get request")

	q := req.URL.Query()
	q.Add("author_id", "A001")

	req.URL.RawQuery = q.Encode()

	resp, err := s.c.Do(req)
	s.Require().NoError(err, "make the actual get request")

	defer resp.Body.Close()

	s.Equal(http.StatusBadRequest, resp.StatusCode)

	var errMsg rest.WebAPIErrorResponse
	err = json.NewDecoder(resp.Body).Decode(&errMsg)
	s.Require().NoError(err, "decode web API error message")
	s.Equal(rest.BadRequestCode, errMsg.Code)
	s.Equal("invalid query params provided", errMsg.Message)
}
