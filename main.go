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
)

const (
	ReadTimeout  int = 5
	WriteTimeout int = 5
)

type health struct {
	Status  string   `json:"status"`
	Message []string `json:"msg"`
}

func main() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("open database: %s", err.Error())
	}

	pgDriver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("create postgresql driver: %s", err.Error())
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations/postgresql", "postgres", pgDriver)
	if err != nil {
		log.Fatalf("make migrations engine: %s", err.Error())
	}

	if err = m.Steps(4); err != nil {
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
