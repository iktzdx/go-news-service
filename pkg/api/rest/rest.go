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
		Payload: createRespPayload(post),
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

	list, err := h.port.List(params)
	if err != nil {
		WrapErrorWithStatus(w, err)

		return
	}

	resp := ListPostsResponse{
		Posts: createRespPayloads(list),
		Total: list.Total,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
}

func (h PrimaryAdapter) Create(w http.ResponseWriter, r *http.Request) {
	var reqBody CreatePostRequest

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		WrapErrorWithStatus(w, err)

		return
	}

	newPost := createReqPayload(reqBody)

	postID, err := h.port.Create(newPost)
	if err != nil {
		WrapErrorWithStatus(w, err)

		return
	}

	resp := CreatePostResponse{
		ID: postID,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
}

func createReqPayload(r CreatePostRequest) posts.Post {
	return posts.Post{
		AuthorID:  r.AuthorID,
		Title:     r.Title,
		Content:   r.Content,
		CreatedAt: r.CreatedAt,
	}
}

func createRespPayload(p posts.Post) Payload {
	return Payload{
		ID:        p.ID,
		AuthorID:  p.AuthorID,
		Title:     p.Title,
		Content:   p.Content,
		CreatedAt: p.CreatedAt,
	}
}

func createRespPayloads(ps posts.Posts) []Payload {
	result := make([]Payload, len(ps.Posts))

	for idx, post := range ps.Posts {
		result[idx] = createRespPayload(post)
	}

	return result
}
