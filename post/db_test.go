package post_test

import (
	"database/sql"
	"gonews/api"
	"gonews/post"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
)

type RepoFindPostByIDSuite struct {
	suite.Suite
	db      *sql.DB
	adapter post.PGSQLSecondaryAdapter
}

func TestRepoFindPostByIDSuite(t *testing.T) {
	suite.Run(t, new(RepoFindPostByIDSuite))
}

func (s *RepoFindPostByIDSuite) SetupTest() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	s.Require().NoError(err, "open database connection")

	s.db = db
	s.adapter = post.NewPGSQLSecondaryAdapter(db)
}

func (s *RepoFindPostByIDSuite) TearDownTest() {
	err := s.db.Close()
	s.Require().NoError(err, "close db connection")
}

func (s *RepoFindPostByIDSuite) TestFindPostThatDoesNotExist() {
	got, err := s.adapter.FindPostByID(12345)
	s.Require().ErrorIs(err, api.ErrPostNotFound)
	s.Zero(got)
}

func (s *RepoFindPostByIDSuite) TestFindPostThatDoesExist() {
	_, err := s.db.Exec(
		"INSERT INTO posts (id, author_id, title, content, created_at) VALUES ($1, $2, $3, $4, $5)",
		42069, 0, "The Future of Sustainable Energy", "The global pursuit of renewable energy sources continues to gain momentum.", 0,
	)
	s.Require().NoError(err, "insert data")

	got, err := s.adapter.FindPostByID(42069)
	s.Require().NoError(err)

	expected := api.Post{
		ID:        42069,
		AuthorID:  0,
		Title:     "The Future of Sustainable Energy",
		Content:   "The global pursuit of renewable energy sources continues to gain momentum.",
		CreatedAt: 0,
	}

	s.Equal(expected, got)
}

func (s *RepoFindPostByIDSuite) TestFindPostUnexpectedError() {
	s.db.Close()

	got, err := s.adapter.FindPostByID(12345)
	s.Require().ErrorIs(err, api.ErrUnexpected)
	s.Zero(got)
}
