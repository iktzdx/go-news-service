package rest_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/iktzdx/skillfactory-gonews/internal/app/posts"
	"github.com/iktzdx/skillfactory-gonews/pkg/api/rest"
)

type CreatePostSuite struct {
	suite.Suite
	req     *http.Request
	resp    *httptest.ResponseRecorder
	port    *rest.MockBoundaryPort
	adapter rest.PrimaryAdapter
}

func TestCreatePostSuite(t *testing.T) {
	suite.Run(t, new(CreatePostSuite))
}

func (s *CreatePostSuite) SetupTest() {
	s.resp = httptest.NewRecorder()

	s.port = new(rest.MockBoundaryPort)
	s.adapter = rest.NewPrimaryAdapter(s.port)
}

func (s *CreatePostSuite) TestCreatePostSucceeded() {
	reqBody := []byte(`{
        "authorId": 0,
        "title": "The Future of Sustainable Energy",
        "content": "The global pursuit of renewable energy sources continues to gain momentum.",
        "createdAt": 1257894000
    }`)

	bReader := bytes.NewReader(reqBody)

	req, err := http.NewRequest(http.MethodPost, "/post", bReader)
	s.Require().NoError(err, "make post request")

	s.req = req

	mockBody := posts.Post{ //nolint:exhaustruct
		AuthorID:  0,
		Title:     "The Future of Sustainable Energy",
		Content:   "The global pursuit of renewable energy sources continues to gain momentum.",
		CreatedAt: 1257894000,
	}

	s.port.On("Create", mockBody).Return(int64(12345), nil)
	s.adapter.Create(s.resp, s.req)
	s.Equal(http.StatusOK, s.resp.Code)
}
