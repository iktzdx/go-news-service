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

type DBPostRetrieverSuite struct {
	suite.Suite
	db *sql.DB
	r  post.DBPostRetriever
}

func TestDBPostRetrieverSuite(t *testing.T) {
	suite.Run(t, new(DBPostRetrieverSuite))
}

func (s *DBPostRetrieverSuite) SetupTest() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	s.Require().NoError(err, "open database connection")

	s.db = db
	s.r = post.NewDBPostRetriever(db)
}

func (s *DBPostRetrieverSuite) TearDownTest() {
	err := s.db.Close()
	s.Require().NoError(err, "close db connection")
}

func (s *DBPostRetrieverSuite) TestRetrievePostThatDoesNotExist() {
	_, err := s.r.FindPostByID(12345)
	s.Require().ErrorIs(err, api.ErrPostNotFound)
}

func (s *DBPostRetrieverSuite) TestRetrievePostThatDoesExist() {
	_, err := s.db.Exec(
		"INSERT INTO posts (id, author_id, title, content, created_at) VALUES ($1, $2, $3, $4, $5)",
		42069, 0, "The Future of Sustainable Energy", "The global pursuit of renewable energy sources continues to gain momentum.", 0,
	)
	s.Require().NoError(err, "insert data")

	retrieved, err := s.r.FindPostByID(42069)
	s.Require().NoError(err)

	post := api.Post{
		ID:        42069,
		AuthorID:  0,
		Title:     "The Future of Sustainable Energy",
		Content:   "The global pursuit of renewable energy sources continues to gain momentum.",
		CreatedAt: 0,
	}

	s.Equal(post, retrieved)
}

func (s *DBPostRetrieverSuite) TestRetrievePostUnexpectedError() {
	s.db.Close()

	_, err := s.r.FindPostByID(12345)
	s.Require().ErrorIs(err, api.ErrUnexpected)
}
