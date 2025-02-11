package controller

import (
	"exampleapp/domain/usecase"
	"fmt"
	"net/http"
)

type HeadersController struct {
	useCase *usecase.GetProductUseCase
}

func NewHeadersController(useCase *usecase.GetProductUseCase) *HeadersController {
	return &HeadersController{useCase}
}

func (c *HeadersController) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}

	p := c.useCase.Handle(usecase.GetProductQuery{2})

	switch {
	case p != nil:
		fmt.Fprintf(w, "db test: %v\n", p.Name)
	default:
		fmt.Fprintf(w, "db test: ---\n")
	}
}
