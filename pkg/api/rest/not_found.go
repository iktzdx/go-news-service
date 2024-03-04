package rest

import (
	"net/http"
)

func NotFound(w http.ResponseWriter, r *http.Request) {
	WrapErrorWithStatus(
		w,
		WebAPIErrorResponse{Code: "004", Message: "no route found"},
		http.StatusNotFound,
	)
}
