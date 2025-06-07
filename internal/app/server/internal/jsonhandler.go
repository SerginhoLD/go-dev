package internal

import (
	"encoding/json"
	"net/http"
)

func jsonResponse(w http.ResponseWriter, v any, code int) {
	w.Header().Del("Content-Length") // @see http.Error
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}

func jsonErrorResponse(w http.ResponseWriter, error string, code int) {
	jsonResponse(w, map[string]string{"message": error}, code)
}
