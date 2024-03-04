package posts

import (
	"strconv"

	"github.com/pkg/errors"

	"github.com/iktzdx/skillfactory-gonews/pkg/storage"
)

func (port BoundaryPort) GetPostByID(id string) (Post, error) {
	postID, err := strconv.Atoi(id)
	if err != nil {
		return Post{}, errors.Wrap(ErrInvalidQueryParam, "parse int")
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
