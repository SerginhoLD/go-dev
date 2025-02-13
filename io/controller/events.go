package controller

import (
	"net/http"
	"time"
)

type ResponseEvent struct {
	Request    *http.Request
	StatusCode int
	Duration   time.Duration
}
