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
	postID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return Post{}, errors.Wrap(ErrInvalidQueryParam, "parse int")
	}

	data, err := port.repo.FindPostByID(postID)
	if err != nil {
		return Post{}, errors.Wrap(err, "get post")
	}

	return FromRepo(data), nil
}

func (port BoundaryPort) List(params QueryParams) (Posts, error) {
	opts, err := mapSearchOpts(params)
	if err != nil {
		return Posts{}, errors.Wrap(ErrInvalidQueryParam, "validate query params")
	}

	bulkData, err := port.repo.List(opts)
	if err != nil {
		return Posts{}, errors.Wrap(err, "list posts")
	}

	return FromRepoBulk(bulkData), nil
}

func mapSearchOpts(params QueryParams) (storage.SearchOpts, error) {
	var opts storage.SearchOpts

	if params.ID != "" {
		postID, err := strconv.ParseInt(params.ID, 10, 64)
		if err != nil {
			return opts, errors.Wrap(err, "parse post id")
		}

		opts.ID = postID
	}

	if params.AuthorID != "" {
		authorID, err := strconv.ParseInt(params.AuthorID, 10, 64)
		if err != nil {
			return opts, errors.Wrap(err, "parse author id")
		}

		opts.AuthorID = authorID
	}

	if params.Limit != "" {
		limit, err := strconv.Atoi(params.Limit)
		if err != nil {
			return opts, errors.Wrap(err, "parse limit")
		}

		opts.Limit = limit
	}

	if params.Offset != "" {
		offset, err := strconv.Atoi(params.Offset)
		if err != nil {
			return opts, errors.Wrap(err, "parse offset")
		}

		opts.Offset = offset
	}

	return opts, nil
}
