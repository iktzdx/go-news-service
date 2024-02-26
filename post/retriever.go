package post

import (
	"gonews/api"

	"github.com/pkg/errors"
)

type PostFinder interface {
	FindPostByID(id int) (api.Post, error)
}

type PostRetriever struct {
	adapter PostFinder
}

func (r PostRetriever) GetPost(id int) (api.Post, error) {
	post, err := r.adapter.FindPostByID(id)
	if err != nil {
		return api.Post{}, errors.Wrap(err, "get post")
	}

	return post, nil
}

func NewPostRetriever(adapter PostFinder) PostRetriever {
	return PostRetriever{adapter}
}
