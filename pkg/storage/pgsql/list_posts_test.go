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

type ListPostsSuite struct {
	suite.Suite
	db      *sql.DB
	adapter pgsql.SecondaryAdapter
}

func TestListPostsSuite(t *testing.T) {
	suite.Run(t, new(ListPostsSuite))
}

func (s *ListPostsSuite) SetupTest() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	s.Require().NoError(err, "open database connection")

	s.db = db
	s.adapter = pgsql.NewSecondaryAdapter(db)
}

func (s *ListPostsSuite) TearDownTest() {
	err := s.db.Close()
	s.Require().NoError(err, "close db connection")
}

func (s *ListPostsSuite) TestListPostsEmptyResult() {
	opts := storage.SearchOpts{
		FilterOpts: storage.FilterOpts{
			ID:       12345,
			AuthorID: 543,
		},
		PaginationOpts: storage.PaginationOpts{
			Limit:  5,
			Offset: 5,
		},
	}

	got, err := s.adapter.List(opts)

	s.Require().NoError(err)
	s.Zero(got)
}

func (s *ListPostsSuite) TestListPostsThatDoExist() {
	//nolint:exhaustruct
	opts := storage.SearchOpts{
		FilterOpts: storage.FilterOpts{
			AuthorID: 1,
		},
		PaginationOpts: storage.PaginationOpts{
			Offset: 1,
		},
	}

	_, err := s.db.Exec(
		"INSERT INTO authors (id, name) VALUES ($1, $2)",
		1, "Test Author",
	)
	s.Require().NoError(err, "insert data to the authors table")

	_, err = s.db.Exec(
		"INSERT INTO posts (id, author_id, title, content, created_at) VALUES ($1, $2, $3, $4, $5)",
		69420, 1, "This test is green #1", "The content for the green test #1", 20120822193532,
	)
	s.Require().NoError(err, "insert data to the posts table")

	_, err = s.db.Exec(
		"INSERT INTO posts (id, author_id, title, content, created_at) VALUES ($1, $2, $3, $4, $5)",
		42420, 1, "This test is green #2", "The content for the green test #2", 20130822193535,
	)
	s.Require().NoError(err, "insert data to the posts table")

	want := storage.BulkData{
		Posts: []storage.Data{
			{
				ID:        42420,
				AuthorID:  1,
				Title:     "This test is green #2",
				Content:   "The content for the green test #2",
				CreatedAt: 20130822193535,
			},
		},
		Total: 1,
	}

	got, err := s.adapter.List(opts)

	s.Require().NoError(err)
	s.Equal(want, got)
}

func (s *ListPostsSuite) TestListPostsUnexpectedError() {
	s.db.Close()

	got, err := s.adapter.List(storage.SearchOpts{}) //nolint:exhaustruct

	s.Require().Error(err)
	s.Require().ErrorContains(err, "database is closed")
	s.Zero(got)
}
