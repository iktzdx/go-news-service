package posts

import (
	"strconv"

	"github.com/pkg/errors"

	"github.com/iktzdx/skillfactory-gonews/internal/app/rest"
	"github.com/iktzdx/skillfactory-gonews/pkg/storage"
)

type BoundaryPort struct {
	repo storage.BoundaryRepoPort
}

func (port BoundaryPort) GetPostByID(id string) (rest.Post, error) {
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

func NewBoundaryPort(repo storage.BoundaryRepoPort) BoundaryPort {
	return BoundaryPort{repo}
}
