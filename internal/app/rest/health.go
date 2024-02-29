package rest

import (
	"encoding/json"
	"net/http"
)

type healthCheckResponse struct {
	Status  string   `json:"status"`
	Message []string `json:"msg"`
}

func HealthCheck(w http.ResponseWriter, _ *http.Request) {
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
}
