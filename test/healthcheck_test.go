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

type EndToEndSuite struct {
	suite.Suite
}

func TestEndToEndSuite(t *testing.T) {
	suite.Run(t, new(EndToEndSuite))
}

func (s *EndToEndSuite) TestHealthCheck() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	c := http.DefaultClient

	req, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:8080/health", nil)
	s.Require().NoError(err)

	req = req.WithContext(ctx)

	r, err := c.Do(req)
	s.Require().NoError(err)

	defer r.Body.Close()

	b, err := io.ReadAll(r.Body)
	s.Require().NoError(err)

	expectedBody := `{"status": "ok", "msg": []}`

	s.Equal(http.StatusOK, r.StatusCode)
	s.JSONEq(expectedBody, string(b))
}
