package post

import (
	"strconv"

	"github.com/pkg/errors"

	"host.local/gonews/api"
)

type PostsBoundaryRepoPort interface {
	FindPostByID(id int) (api.Post, error)
}

type PostsBoundaryPort struct {
	repo PostsBoundaryRepoPort
}

func (port PostsBoundaryPort) GetPost(id string) (api.Post, error) {
	postID, err := strconv.Atoi(id)
	if err != nil {
		return api.Post{}, errors.Wrap(api.ErrInvalidPostID, "parse int")
	}

	post, err := port.repo.FindPostByID(postID)
	if err != nil {
		return api.Post{}, errors.Wrap(err, "get post")
	}

	return post, nil
}

func NewPostsBoundaryPort(repo PostsBoundaryRepoPort) PostsBoundaryPort {
	return PostsBoundaryPort{repo}
}
