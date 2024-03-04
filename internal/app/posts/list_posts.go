package posts

import (
	"strconv"

	"github.com/pkg/errors"

	"github.com/iktzdx/skillfactory-gonews/pkg/storage"
)

func (port BoundaryPort) List(params QueryParams) (Posts, error) {
	var posts Posts

	_, err := validateQueryParams(params)
	if err != nil {
		return posts, errors.Wrap(ErrInvalidQueryParam, "validate query params")
	}

	return posts, nil
}

func validateQueryParams(params QueryParams) (storage.SearchOpts, error) {
	var opts storage.SearchOpts

	if params.ID != "" {
		postID, err := strconv.Atoi(params.ID)
		if err != nil {
			return opts, errors.Wrap(err, "parse post id")
		}

		opts.ID = postID
	}

	if params.AuthorID != "" {
		authorID, err := strconv.Atoi(params.AuthorID)
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
