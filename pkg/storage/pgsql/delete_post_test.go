//go:build integration

package pgsql_test

import (
	"database/sql"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"

	"github.com/iktzdx/skillfactory-gonews/pkg/storage"
	"github.com/iktzdx/skillfactory-gonews/pkg/storage/pgsql"
)

type PQDeletePostSuite struct {
	suite.Suite
	db      *sql.DB
	adapter pgsql.SecondaryAdapter
}

func TestPQDeletePostSuite(t *testing.T) {
	suite.Run(t, new(PQDeletePostSuite))
}

func (s *PQDeletePostSuite) SetupTest() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	s.Require().NoError(err, "open database connection")

	s.db = db
	s.adapter = pgsql.NewSecondaryAdapter(db)
}

func (s *PQDeletePostSuite) TearDownTest() {
	err := s.db.Close()
	s.Require().NoError(err, "close db connection")
}

func (s *PQDeletePostSuite) TestPQDeletePostThatDoesNotExist() {
	got, err := s.adapter.Delete(12345)

	s.Require().ErrorIs(err, storage.ErrNoDataFound)
	s.Require().EqualValues(-1, got)
}

func (s *PQDeletePostSuite) TestPQDeletePostThatDoesExist() {
	now := time.Now().UTC().Unix()

	data := storage.Data{ //nolint:exhaustruct
		AuthorID:  0,
		Title:     "Test post",
		Content:   "Pls delete me!",
		CreatedAt: now,
	}

	createdID, err := s.adapter.Create(data)
	s.Require().NoError(err)

	got, err := s.adapter.Delete(createdID)

	s.Require().NoError(err)
	s.Require().EqualValues(1, got)
}

func (s *PQDeletePostSuite) TestPQDeletePostUnexpectedError() {
	s.db.Close()

	got, err := s.adapter.Delete(54321) //nolint:exhaustruct

	s.Require().Error(err)
	s.Require().ErrorContains(err, "database is closed")
	s.Require().EqualValues(-1, got)
}
