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

type PostRetrieverSuite struct {
	suite.Suite
	db *sql.DB
	r  post.PostRetriever
}

func TestPostRetrieverSuite(t *testing.T) {
	suite.Run(t, new(PostRetrieverSuite))
}

func (s *PostRetrieverSuite) SetupTest() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	s.Require().NoError(err, "open database connection")

	s.db = db
	s.r = post.NewPostRetriever(db)
}

func (s *PostRetrieverSuite) TearDownTest() {
	err := s.db.Close()
	s.Require().NoError(err, "close db connection")
}

func (s *PostRetrieverSuite) TestRetrievePostThatDoesNotExist() {
	_, err := s.r.GetPost(12345)
	s.Require().ErrorIs(err, api.ErrPostNotFound)
}

func (s *PostRetrieverSuite) TestRetrievePostThatDoesExist() {
	_, err := s.db.Exec(
		"INSERT INTO posts (id, author_id, title, content, created_at) VALUES ($1, $2, $3, $4, $5)",
		42069, 0, "The Future of Sustainable Energy", "The global pursuit of renewable energy sources continues to gain momentum.", 0,
	)
	s.Require().NoError(err, "insert data")

	retrieved, err := s.r.GetPost(42069)
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

func (s *PostRetrieverSuite) TestRetrievePostUnexpectedError() {
	s.db.Close()

	_, err := s.r.GetPost(12345)
	s.Require().ErrorIs(err, api.ErrUnexpected)
}
