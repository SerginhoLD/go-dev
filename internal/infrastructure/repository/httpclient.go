package repository

import (
	"net/http"
	"os"
	"strconv"
	"time"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func NewHttpClient() HttpClient {
	timeout, _ := strconv.ParseUint(os.Getenv("OBJECT_HTTP_TIMEOUT"), 10, 64)

	return &http.Client{
		Timeout: time.Second * time.Duration(timeout),
	}
}

func NewHttpRequest(method string, url string) *http.Request {
	req, _ := http.NewRequest(method, url, nil)
	req.Header.Set("User-Agent", os.Getenv("OBJECT_USER_AGENT"))

	return req
}
