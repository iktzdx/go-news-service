package rest

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

type PrimaryAdapter struct {
	port BoundaryPort
}

func NewPrimaryAdapter(port BoundaryPort) PrimaryAdapter {
	return PrimaryAdapter{port}
}

func (h PrimaryAdapter) GetPostByID(w http.ResponseWriter, r *http.Request) {
	postID := mux.Vars(r)["id"]

	post, err := h.port.GetPostByID(postID)
	if err != nil {
		var status int

		var errMsg WebAPIError

		switch {
		case errors.Is(err, ErrPostNotFound):
			status = http.StatusNotFound
			errMsg = WebAPIError{Code: "001", Message: "no post with id " + postID}
		case errors.Is(err, ErrInvalidPostID):
			status = http.StatusBadRequest
			errMsg = WebAPIError{Code: "003", Message: "invalid post id provided"}
		default:
			status = http.StatusInternalServerError
			errMsg = WebAPIError{"002", "unexpected error attempting to get post"}
		}

		WrapErrorWithStatus(w, errMsg, status)

		return
	}

	WrapOK(w, post)
}
