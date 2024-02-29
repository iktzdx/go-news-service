package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/lib/pq"

	"github.com/iktzdx/skillfactory-gonews/internal/app/posts"
	"github.com/iktzdx/skillfactory-gonews/pkg/api/rest"
	"github.com/iktzdx/skillfactory-gonews/pkg/storage/pgsql"
)

const (
	ReadTimeout    time.Duration = 5 * time.Second
	WriteTimeout   time.Duration = 5 * time.Second
	MigrationSteps int           = 4
)

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

	repo := pgsql.NewSecondaryAdapter(db)
	boundaryPort := posts.NewBoundaryPort(repo)
	postsHandler := rest.NewPrimaryAdapter(boundaryPort)

	routes := rest.CreateRoutes(postsHandler)

	s := http.Server{ //nolint:exhaustruct
		Addr:         ":8080",
		Handler:      routes,
		ReadTimeout:  ReadTimeout,
		WriteTimeout: WriteTimeout,
	}

	if err := s.ListenAndServe(); err != nil {
		log.Fatalf("listen on %s: %s", s.Addr, err.Error())
	}
}
