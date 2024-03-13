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

type PQUpdatePostSuite struct {
	suite.Suite
	db      *sql.DB
	adapter pgsql.SecondaryAdapter
}

func TestPQUpdatePostSuite(t *testing.T) {
	suite.Run(t, new(PQUpdatePostSuite))
}

func (s *PQUpdatePostSuite) SetupTest() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	s.Require().NoError(err, "open database connection")

	s.db = db
	s.adapter = pgsql.NewSecondaryAdapter(db)
}

func (s *PQUpdatePostSuite) TearDownTest() {
	err := s.db.Close()
	s.Require().NoError(err, "close db connection")
}

func (s *PQUpdatePostSuite) TestPQUpdatePostThatDoesNotExist() {
	change := storage.Data{ //nolint:exhaustruct
		ID:       10293,
		AuthorID: 0,
		Title:    "Non-existent post",
		Content:  "This post does not exist.",
	}

	got, err := s.adapter.Update(change)

	s.Require().ErrorIs(err, storage.ErrNoDataFound)
	s.Require().EqualValues(-1, got)
}

func (s *PQUpdatePostSuite) TestPQUpdatePostThatDoesExist() {
	now := time.Now().UTC().Unix()

	data := storage.Data{ //nolint:exhaustruct
		AuthorID:  1,
		Title:     "New post",
		Content:   "This is a new post for testing.",
		CreatedAt: now,
	}

	createdID, err := s.adapter.Create(data)
	s.Require().NoError(err)

	want := storage.Data{
		ID:        createdID,
		AuthorID:  0,
		Title:     "Updated post title",
		Content:   "This post was updated.",
		CreatedAt: now,
	}

	affected, err := s.adapter.Update(want)
	s.Require().NoError(err)
	s.Require().Greater(affected, int64(0), "affected rows")

	got, err := s.adapter.FindPostByID(createdID)
	s.Require().NoError(err)
	s.Require().Equal(want, got)
}

func (s *PQUpdatePostSuite) TestPQUpdatePostNonExistentAuthorID() {
	now := time.Now().UTC().Unix()

	data := storage.Data{ //nolint:exhaustruct
		AuthorID:  1,
		Title:     "New post",
		Content:   "This is a new post for testing.",
		CreatedAt: now,
	}

	createdID, err := s.adapter.Create(data)
	s.Require().NoError(err)

	change := storage.Data{
		ID:        createdID,
		AuthorID:  42069,
		Title:     "Updated post title",
		Content:   "This post was updated.",
		CreatedAt: now,
	}

	got, err := s.adapter.Update(change)

	s.Require().Error(err)
	s.Require().ErrorContains(err, "violates foreign key constraint \"posts_author_id_fkey\"")
	s.Require().EqualValues(-1, got)
}
