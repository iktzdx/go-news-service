//go:build e2e

package e2e_test

import (
	"context"
	"io"
	"net/http"
	"testing"
	"time"

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
