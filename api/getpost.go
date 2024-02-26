package api

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"

	"github.com/gorilla/mux"
)

type PostsBoundaryPort interface {
	GetPost(id string) (Post, error)
}

var (
	ErrPostNotFound  = errors.New("post not found")
	ErrInvalidPostID = errors.New("invalid post id")
	ErrUnexpected    = errors.New("unexpected error")
)

type Post struct {
	ID        int    `json:"id"`
	AuthorID  int    `json:"authorId"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	CreatedAt int    `json:"createdAt"`
}

type WebAPIError struct {
	Code    string `json:"code"`
	Message string `json:"msg"`
}

type RESTPrimaryAdapter struct {
	port PostsBoundaryPort
}

func (h RESTPrimaryAdapter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	postID := mux.Vars(r)["id"]

	post, err := h.port.GetPost(postID)
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

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)

		if err := json.NewEncoder(w).Encode(errMsg); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(post); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
}

func NewRESTPrimaryAdapter(port PostsBoundaryPort) RESTPrimaryAdapter {
	return RESTPrimaryAdapter{port}
}
