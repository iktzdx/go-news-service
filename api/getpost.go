package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

type PostRetriever interface {
	GetPost(id int) (Post, error)
}

var ErrPostNotFound = errors.New("post not found")

type Post struct {
	ID        int    `json:"id"`
	AuthorID  int    `json:"authorId"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	CreatedAt int    `json:"createdAt"`
}

type errWebAPI struct {
	Code    string `json:"code"`
	Message string `json:"msg"`
}

type GetPosthandler struct{}

func (h *GetPosthandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	postID := mux.Vars(r)["id"]

	errPostNotFound := errWebAPI{
		Code:    "001",
		Message: "no post with id " + postID,
	}

	body, err := json.Marshal(errPostNotFound)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)

	if _, err = w.Write(body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
}

func NewGetPostHandler(r PostRetriever) *GetPosthandler {
	return new(GetPosthandler)
}
