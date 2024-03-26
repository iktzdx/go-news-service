//go:build e2e

package test_test

import (
	"bytes"
	"net/http"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
)

type CreatePostSuite struct {
	suite.Suite
	c *http.Client
}

func TestCreatePostSuite(t *testing.T) {
	suite.Run(t, new(CreatePostSuite))
}

func (s *CreatePostSuite) SetupTest() {
	s.c = http.DefaultClient
}

func (s *CreatePostSuite) TestCreatePostSucceeded() {
	reqBody := []byte(`{
        "authorId": 0,
        "title": "The Future of Sustainable Energy",
        "content": "The global pursuit of renewable energy sources continues to gain momentum.",
        "createdAt": 1257894000
    }`)

	bReader := bytes.NewReader(reqBody)

	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:8080/post", bReader) //nolint:noctx
	s.Require().NoError(err, "make post request")

	resp, err := s.c.Do(req)
	s.Require().NoError(err, "send post request")

	defer resp.Body.Close()

	s.Equal(http.StatusOK, resp.StatusCode)
}
