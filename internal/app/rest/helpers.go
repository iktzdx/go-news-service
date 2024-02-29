package rest

import (
	"encoding/json"
	"net/http"
)

func WrapErrorWithStatus(w http.ResponseWriter, errMsg WebAPIError, status int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(errMsg); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
}

func WrapOK(w http.ResponseWriter, post Post) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(post); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
}
