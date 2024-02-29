package api

import (
	"github.com/gorilla/mux"

	"github.com/iktzdx/skillfactory-gonews/internal/app/rest"
)

func CreateRoutes(adapter rest.PrimaryAdapter) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/health", rest.HealthCheck).Methods("GET")

	// r.HandleFunc("/post", adapter.Create).Methods("POST")
	r.HandleFunc("/post/{id}", adapter.GetPostByID).Methods("GET")
	// r.HandleFunc("/posts", adapter.List).Methods("GET")
	// r.HandleFunc("/post/{id}", adapter.Update).Methods("PUT")
	// r.HandleFunc("/post/{id}", adapter.Delete).Methods("DELETE")

	r.NotFoundHandler = r.NewRoute().HandlerFunc(rest.NotFound).GetHandler()

	return r
}
