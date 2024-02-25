package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type PostRetriever interface {
	GetPost(id int) (Post, error)
}

var (
	ErrPostNotFound = errors.New("post not found")
	ErrUnexpected   = errors.New("unexpected error")
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

type GetPostHandler struct {
	r PostRetriever
}

func (h *GetPostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)["id"]

	postID, err := strconv.Atoi(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	post, err := h.r.GetPost(postID)
	if err != nil {
		var status int

		var errMsg WebAPIError

		switch {
		case errors.Is(err, ErrPostNotFound):
			status = http.StatusNotFound
			errMsg = WebAPIError{Code: "001", Message: "no post with id " + v}
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

func NewGetPostHandler(r PostRetriever) *GetPostHandler {
	return &GetPostHandler{r}
}
