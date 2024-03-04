//go:build e2e

package test_test

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"

	"github.com/iktzdx/skillfactory-gonews/pkg/api/rest"
)

type GetPostByIDSuite struct {
	suite.Suite
	c *http.Client
}

func TestGetPostByIDSuite(t *testing.T) {
	suite.Run(t, new(GetPostByIDSuite))
}

func (s *GetPostByIDSuite) SetupTest() {
	s.c = http.DefaultClient
}

func (s *GetPostByIDSuite) TestGetPostThatDoesNotExist() {
	req, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:8080/post/12345", nil)
	s.Require().NoError(err)

	resp, err := s.c.Do(req)
	s.Require().NoError(err)

	defer resp.Body.Close()

	s.Equal(http.StatusNotFound, resp.StatusCode)

	var errMsg rest.WebAPIErrorResponse
	err = json.NewDecoder(resp.Body).Decode(&errMsg)
	s.Require().NoError(err, "decode web API error message")
	s.Equal("001", errMsg.Code)
	s.Equal("no post with id 12345", errMsg.Message)
}

func (s *GetPostByIDSuite) TestGetPostThatDoesExist() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	s.Require().NoError(err, "open database connection")

	res, err := db.Exec(
		"INSERT INTO posts (id, author_id, title, content, created_at) VALUES ($1, $2, $3, $4, $5)",
		54321, 0, "The Future of Sustainable Energy", "The global pursuit of renewable energy sources continues to gain momentum.", 0,
	)
	s.Require().NoError(err, "insert data")

	affected, err := res.RowsAffected()
	s.Require().NoError(err, "get the number of rows affected")
	s.Require().EqualValues(1, affected, "rows affected")

	req, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:8080/post/54321", nil) //nolint:noctx
	s.Require().NoError(err, "make get request")

	resp, err := s.c.Do(req)
	s.Require().NoError(err, "send get request")

	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	s.Require().NoError(err, "read response body")

	expectedBody := `{
        "id": 54321,
        "authorId": 0,
        "title": "The Future of Sustainable Energy",
        "content": "The global pursuit of renewable energy sources continues to gain momentum.",
        "createdAt": 0
    }`

	s.Equal(http.StatusOK, resp.StatusCode)
	s.JSONEq(expectedBody, string(b))
}

func (s *GetPostByIDSuite) TestGetPostWithInvalidID() {
	req, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:8080/post/12C45", nil) //nolint:noctx
	s.Require().NoError(err)

	resp, err := s.c.Do(req)
	s.Require().NoError(err)

	defer resp.Body.Close()

	s.Equal(http.StatusBadRequest, resp.StatusCode)

	var errMsg rest.WebAPIErrorResponse
	err = json.NewDecoder(resp.Body).Decode(&errMsg)
	s.Require().NoError(err, "decode web API error message")
	s.Equal("003", errMsg.Code)
	s.Equal("invalid post id provided", errMsg.Message)
}
