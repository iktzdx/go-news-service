package rest

import (
	"net/http"
)

func NotFound(w http.ResponseWriter, r *http.Request) {
	WrapErrorWithStatus(
		w,
		WebAPIErrorResponse{Code: RouteNotFoundCode, Message: "no route found"},
		http.StatusNotFound,
	)
}
