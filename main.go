package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"host.local/gonews/api"
	"host.local/gonews/post"
)

const (
	ReadTimeout    time.Duration = 5 * time.Second
	WriteTimeout   time.Duration = 5 * time.Second
	MigrationSteps int           = 4
)

type healthCheckResponse struct {
	Status  string   `json:"status"`
	Message []string `json:"msg"`
}

type Post struct {
	ID        int    `json:"id"`
	AuthorID  int    `json:"authorId"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	CreatedAt int    `json:"createdAt"`
}

func main() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("open database: %s", err.Error())
	}

	var cfg postgres.Config

	pgDriver, err := postgres.WithInstance(db, &cfg)
	if err != nil {
		log.Fatalf("create postgresql driver: %s", err.Error())
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations/postgresql", "postgres", pgDriver)
	if err != nil {
		log.Fatalf("make migrations engine: %s", err.Error())
	}

	if err = m.Steps(MigrationSteps); err != nil {
		log.Fatalf("migrate up: %s", err.Error())
	}

	r := mux.NewRouter()

	repo := post.NewPGSQLSecondaryAdapter(db)
	boundaryPort := post.NewPostsBoundaryPort(repo)

	r.Handle("/post/{id}", api.NewRESTPrimaryAdapter(boundaryPort))

	r.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		resp, err := json.Marshal(healthCheckResponse{
			Status:  "OK",
			Message: []string{},
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if _, err = w.Write(resp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}
	})

	s := http.Server{ //nolint:exhaustruct
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  ReadTimeout,
		WriteTimeout: WriteTimeout,
	}

	if err := s.ListenAndServe(); err != nil {
		log.Fatalf("listen on %s: %s", s.Addr, err.Error())
	}
}
