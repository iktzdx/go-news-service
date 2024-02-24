//go:build e2e

package e2e_test

import (
	"context"
	"database/sql"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
)

type GetPostByIDSuite struct {
	suite.Suite
}

func TestGetPostByIDSuite(t *testing.T) {
	suite.Run(t, new(GetPostByIDSuite))
}

func (s *GetPostByIDSuite) TestGetPostThatDoesNotExist() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	c := http.DefaultClient

	req, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:8080/post/12345", nil)
	s.Require().NoError(err)

	req = req.WithContext(ctx)

	r, err := c.Do(req)
	s.Require().NoError(err)

	defer r.Body.Close()

	b, err := io.ReadAll(r.Body)
	s.Require().NoError(err)

	expectedBody := `{"code": "001", "msg": "no post with id 12345"}`

	s.Equal(http.StatusNotFound, r.StatusCode)
	s.JSONEq(expectedBody, string(b))
}

func (s *GetPostByIDSuite) TestGetPostThatDoesExist() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

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

	c := http.DefaultClient

	req, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:8080/post/54321", nil)
	s.Require().NoError(err, "make get request")

	req = req.WithContext(ctx)

	r, err := c.Do(req)
	s.Require().NoError(err, "send get request")

	defer r.Body.Close()

	b, err := io.ReadAll(r.Body)
	s.Require().NoError(err, "read response body")

	expectedBody := `{
        "id": 54321,
        "authorId": 0,
        "title": "The Future of Sustainable Energy",
        "content": "The global pursuit of renewable energy sources continues to gain momentum.",
        "createdAt": 0
    }`

	s.Equal(http.StatusOK, r.StatusCode)
	s.JSONEq(expectedBody, string(b))
}
