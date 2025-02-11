package controller

import (
	"exampleapp/domain/repository"
	"fmt"
	"net/http"
)

type HeadersController struct {
	repository repository.ProductRepository
}

func NewHeadersController(repository repository.ProductRepository) *HeadersController {
	return &HeadersController{repository}
}

func (c *HeadersController) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}

	p := c.repository.Find(11)

	switch {
	case p != nil:
		fmt.Fprintf(w, "db test: %v\n", p.Name)
	default:
		fmt.Fprintf(w, "db test: ---\n")
	}
}
