package rest

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.com/iktzdx/skillfactory-gonews/internal/app/posts"
	"github.com/iktzdx/skillfactory-gonews/pkg/api"
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
		WrapErrorWithStatus(w, err)

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

func (h PrimaryAdapter) List(w http.ResponseWriter, r *http.Request) {
	vars := r.URL.Query()

	params := posts.QueryParams{
		FiltersParams: posts.FiltersParams{
			ID:       strings.Trim(vars.Get("id"), "\""),
			AuthorID: strings.Trim(vars.Get("author_id"), "\""),
		},
		PaginationParams: posts.PaginationParams{
			Limit:  strings.Trim(vars.Get("limit"), "\""),
			Offset: strings.Trim(vars.Get("offset"), "\""),
		},
	}

	result, err := h.port.List(params)
	if err != nil {
		WrapErrorWithStatus(w, err)

		return
	}

	posts := make([]Payload, len(result.Posts))
	for idx, post := range result.Posts {
		posts[idx] = Payload{
			ID:        post.ID,
			AuthorID:  post.AuthorID,
			Title:     post.Title,
			Content:   post.Content,
			CreatedAt: post.CreatedAt,
		}
	}

	resp := ListPostsResponse{
		Posts: posts,
		Total: result.Total,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
}
