//go:build integration

package pgsql_test

import (
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/iktzdx/skillfactory-gonews/pkg/storage"
	"github.com/iktzdx/skillfactory-gonews/pkg/storage/pgsql"
)

type CreatePostSuite struct {
	suite.Suite
	db      *sql.DB
	adapter pgsql.SecondaryAdapter
}

func TestCreatePostSuite(t *testing.T) {
	suite.Run(t, new(CreatePostSuite))
}

func (s *CreatePostSuite) SetupTest() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	s.Require().NoError(err, "open database connection")

	s.db = db
	s.adapter = pgsql.NewSecondaryAdapter(db)
}

func (s *CreatePostSuite) TearDownTest() {
	err := s.db.Close()
	s.Require().NoError(err, "close db connection")
}

func (s *CreatePostSuite) TestCreatePostSucceed() {
	now := time.Now().UTC().Unix()

	want := storage.Data{ //nolint:exhaustruct
		AuthorID:  0,
		Title:     "Test post title",
		Content:   "Just a simple text data.",
		CreatedAt: now,
	}

	createdID, err := s.adapter.Create(want)
	s.Require().NoError(err)
	s.Require().NotZero(createdID)

	want.ID = createdID

	got, err := s.adapter.FindPostByID(want.ID)
	s.Require().NoError(err)
	s.Equal(want, got)
}

func (s *CreatePostSuite) TestCreatePostFailed() {
	s.db.Close()

	got, err := s.adapter.Create(storage.Data{}) //nolint:exhaustruct

	s.Require().Error(err)
	s.Require().ErrorContains(err, "database is closed")
	s.Require().EqualValues(-1, got)
}
