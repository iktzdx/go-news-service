package rest

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"

	"github.com/iktzdx/skillfactory-gonews/internal/app/posts"
	"github.com/iktzdx/skillfactory-gonews/pkg/storage"
)

func WrapErrorWithStatus(w http.ResponseWriter, err error) {
	var status int

	var errMsg WebAPIErrorResponse

	switch {
	case errors.Is(err, posts.ErrInvalidQueryParam):
		status = http.StatusBadRequest
		errMsg = WebAPIErrorResponse{Code: BadRequestCode, Message: "invalid query params provided"}
	case errors.Is(err, storage.ErrNoDataFound):
		status = http.StatusNotFound
		errMsg = WebAPIErrorResponse{Code: PostNotFoundCode, Message: "no post found with id provided"}
	case errors.Is(err, ErrNoRouteFound):
		status = http.StatusNotFound
		errMsg = WebAPIErrorResponse{Code: RouteNotFoundCode, Message: "no route found"}
	default:
		status = http.StatusInternalServerError
		errMsg = WebAPIErrorResponse{Code: UnexpectedCode, Message: "service returned unexpected error"}
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(errMsg); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
}
