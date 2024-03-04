package posts

import (
	"strconv"

	"github.com/pkg/errors"

	"github.com/iktzdx/skillfactory-gonews/pkg/storage"
)

type BoundaryPort struct {
	repo storage.BoundaryRepoPort
}

func NewBoundaryPort(repo storage.BoundaryRepoPort) BoundaryPort {
	return BoundaryPort{repo}
}

func (port BoundaryPort) GetPostByID(id string) (Post, error) {
	postID, err := strconv.Atoi(id)
	if err != nil {
		return Post{}, errors.Wrap(ErrInvalidPostID, "parse int")
	}

	data, err := port.repo.FindPostByID(postID)
	if err != nil {
		return Post{}, errors.Wrap(err, "get post")
	}

	return FromRepo(data), nil
}

func FromRepo(data storage.Data) Post {
	return Post{
		ID:        data.ID,
		AuthorID:  data.AuthorID,
		Title:     data.Title,
		Content:   data.Content,
		CreatedAt: data.CreatedAt,
	}
}
