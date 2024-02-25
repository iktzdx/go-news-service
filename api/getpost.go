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

type GetPosthandler struct {
	r PostRetriever
}

func (h *GetPosthandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)["id"]

	postID, err := strconv.Atoi(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	post, err := h.r.GetPost(postID)
	if err != nil {
		errPostNotFound := errWebAPI{
			Code:    "001",
			Message: "no post with id " + v,
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

		return
	}

	body, err := json.Marshal(post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if _, err = w.Write(body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
}

func NewGetPostHandler(r PostRetriever) *GetPosthandler {
	return &GetPosthandler{r}
}
