package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
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
		log.Fatal(err)
	}
}
