package rest

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"github.com/iktzdx/skillfactory-gonews/internal/app/posts"
	"github.com/iktzdx/skillfactory-gonews/pkg/api"
	"github.com/iktzdx/skillfactory-gonews/pkg/storage"
)

type PrimaryAdapter struct {
	port api.BoundaryPort
}

func NewPrimaryAdapter(port api.BoundaryPort) PrimaryAdapter {
	return PrimaryAdapter{port}
}

func (h PrimaryAdapter) GetPostByID(w http.ResponseWriter, r *http.Request) {
	postID := mux.Vars(r)["id"]

	post, err := h.port.GetPostByID(postID)
	if err != nil {
		var status int

		var errMsg WebAPIErrorResponse

		switch {
		case errors.Is(err, storage.ErrNoDataFound):
			status = http.StatusNotFound
			errMsg = WebAPIErrorResponse{Code: "001", Message: "no post with id " + postID}
		case errors.Is(err, posts.ErrInvalidPostID):
			status = http.StatusBadRequest
			errMsg = WebAPIErrorResponse{Code: "003", Message: "invalid post id provided"}
		default:
			status = http.StatusInternalServerError
			errMsg = WebAPIErrorResponse{"002", "unexpected error attempting to get post"}
		}

		WrapErrorWithStatus(w, errMsg, status)

		return
	}

	resp := GetPostByIDResponse{
		Payload: Payload{
			ID:        post.ID,
			AuthorID:  post.AuthorID,
			Title:     post.Title,
			Content:   post.Content,
			CreatedAt: post.CreatedAt,
		},
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
}
