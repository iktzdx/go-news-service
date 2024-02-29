package posts

import (
	"strconv"

	"github.com/pkg/errors"

	"github.com/iktzdx/skillfactory-gonews/internal/app/rest"
)

type BoundaryRepoPort interface {
	FindPostByID(id int) (rest.Post, error)
}

type BoundaryPort struct {
	repo BoundaryRepoPort
}

func (port BoundaryPort) GetPost(id string) (rest.Post, error) {
	postID, err := strconv.Atoi(id)
	if err != nil {
		return rest.Post{}, errors.Wrap(rest.ErrInvalidPostID, "parse int")
	}

	post, err := port.repo.FindPostByID(postID)
	if err != nil {
		return rest.Post{}, errors.Wrap(err, "get post")
	}

	return post, nil
}

func NewBoundaryPort(repo BoundaryRepoPort) BoundaryPort {
	return BoundaryPort{repo}
}
