package api

import (
	"github.com/gorilla/mux"

	"github.com/iktzdx/skillfactory-gonews/internal/app/rest"
)

func CreateRoutes(adapter rest.PrimaryAdapter) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/health", rest.HealthCheck).Methods("GET")

	r.HandleFunc("/post/{id}", adapter.ServeHTTP).Methods("GET")

	// r.HandleFunc("/users/create", usersHandler.Create).Methods("POST")
	// r.HandleFunc("/users/list", usersHandler.List).Methods("GET")

	r.NotFoundHandler = r.NewRoute().HandlerFunc(rest.NotFound).GetHandler()

	return r
}
