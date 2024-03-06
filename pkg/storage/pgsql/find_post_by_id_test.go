//go:build integration

package pgsql_test

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"

	"github.com/iktzdx/skillfactory-gonews/pkg/storage"
	"github.com/iktzdx/skillfactory-gonews/pkg/storage/pgsql"
)

type FindPostByIDSuite struct {
	suite.Suite
	db      *sql.DB
	adapter pgsql.SecondaryAdapter
}

func TestFindPostByIDSuite(t *testing.T) {
	suite.Run(t, new(FindPostByIDSuite))
}

func (s *FindPostByIDSuite) SetupTest() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	s.Require().NoError(err, "open database connection")

	s.db = db
	s.adapter = pgsql.NewSecondaryAdapter(db)
}

func (s *FindPostByIDSuite) TearDownTest() {
	err := s.db.Close()
	s.Require().NoError(err, "close db connection")
}

func (s *FindPostByIDSuite) TestFindPostThatDoesNotExist() {
	got, err := s.adapter.FindPostByID(12345)
	s.Require().ErrorIs(err, storage.ErrNoDataFound)
	s.Zero(got)
}

func (s *FindPostByIDSuite) TestFindPostThatDoesExist() {
	_, err := s.db.Exec(
		"INSERT INTO posts (id, author_id, title, content, created_at) VALUES ($1, $2, $3, $4, $5)",
		42069, 0, "The Future of Sustainable Energy", "The global pursuit of renewable energy sources continues to gain momentum.", 0,
	)
	s.Require().NoError(err, "insert data")

	got, err := s.adapter.FindPostByID(42069)
	s.Require().NoError(err)

	want := storage.Data{
		ID:        42069,
		AuthorID:  0,
		Title:     "The Future of Sustainable Energy",
		Content:   "The global pursuit of renewable energy sources continues to gain momentum.",
		CreatedAt: 0,
	}

	s.Equal(want, got)
}

func (s *FindPostByIDSuite) TestFindPostUnexpectedError() {
	s.db.Close()

	got, err := s.adapter.FindPostByID(12345)

	s.Require().Error(err)
	s.Require().ErrorContains(err, "database is closed")
	s.Zero(got)
}
