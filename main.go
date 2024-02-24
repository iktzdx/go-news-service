package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	ReadTimeout    int = 5
	WriteTimeout   int = 5
	MigrationSteps int = 4
)

type health struct {
	Status  string   `json:"status"`
	Message []string `json:"msg"`
}

type errWebAPI struct {
	Code    string `json:"code"`
	Message string `json:"msg"`
}

type Post struct {
	ID        int    `json:"id"`
	AuthorID  int    `json:"authorId"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	CreatedAt int    `json:"createdAt"`
}

func main() { //nolint:funlen,cyclop
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

	r.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		resp := health{
			Status:  "ok",
			Message: []string{},
		}

		b, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if _, err = w.Write(b); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}
	})

	r.HandleFunc("/post/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]

		postID, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		query := "SELECT * FROM posts WHERE id = $1"
		row := db.QueryRow(query, postID)

		var post Post
		if err := row.Scan(&post.ID, &post.AuthorID, &post.Title, &post.Content, &post.CreatedAt); err != nil {
			errPostNotFound := errWebAPI{
				Code:    "001",
				Message: "no post with id " + id,
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
	})

	s := http.Server{ //nolint:exhaustruct
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  time.Duration(ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(ReadTimeout) * time.Second,
	}

	if err := s.ListenAndServe(); err != nil {
		log.Fatalf("listen on %s: %s", s.Addr, err.Error())
	}
}
