package controller

import (
	"encoding/json"
	"net/http"
)

func NotFoundHandler(w http.ResponseWriter, req *http.Request) {
	HttpJsonError(w, "Not Found", http.StatusNotFound)
}

func HttpJsonError(w http.ResponseWriter, error string, code int) {
	// @see http.Error
	h := w.Header()

	// Delete the Content-Length header, which might be for some other content.
	// Assuming the error string fits in the writer's buffer, we'll figure
	// out the correct Content-Length for it later.
	//
	// We don't delete Content-Encoding, because some middleware sets
	// Content-Encoding: gzip and wraps the ResponseWriter to compress on-the-fly.
	// See https://go.dev/issue/66343.
	h.Del("Content-Length")

	// There might be content type already set, but we reset it to
	// text/plain for the error message.
	h.Set("Content-Type", "application/json")
	h.Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"message": error})
}
