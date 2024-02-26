package e2e_test

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"
)

type EndToEndSuite struct {
	suite.Suite
}

func TestEndToEndSuite(t *testing.T) {
	suite.Run(t, new(EndToEndSuite))
}

func (s *EndToEndSuite) TestHealthCheck() {
	c := http.DefaultClient

	req, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:8080/health", nil) //nolint:noctx
	s.Require().NoError(err)

	r, err := c.Do(req)
	s.Require().NoError(err)

	defer r.Body.Close()

	b, err := io.ReadAll(r.Body)
	s.Require().NoError(err)

	expectedBody := `{"status": "OK", "msg": []}`

	s.Equal(http.StatusOK, r.StatusCode)
	s.JSONEq(expectedBody, string(b))
}
