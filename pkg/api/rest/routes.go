package rest

import (
	"github.com/gorilla/mux"
)

func CreateRoutes(adapter PrimaryAdapter) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/health", HealthCheck).Methods("GET")

	// r.HandleFunc("/post", adapter.Create).Methods("POST")
	r.HandleFunc("/post/{id}", adapter.GetPostByID).Methods("GET")
	// r.HandleFunc("/posts", adapter.List).Methods("GET")
	// r.HandleFunc("/post/{id}", adapter.Update).Methods("PUT")
	// r.HandleFunc("/post/{id}", adapter.Delete).Methods("DELETE")

	r.NotFoundHandler = r.NewRoute().HandlerFunc(NotFound).GetHandler()

	return r
}
