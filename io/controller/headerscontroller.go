package controller

import (
	"encoding/json"
	"exampleapp/domain/usecase"
	"net/http"
)

type HeadersController struct {
	useCase *usecase.GetProductUseCase
}

func NewHeadersController(useCase *usecase.GetProductUseCase) *HeadersController {
	return &HeadersController{useCase}
}

func (c *HeadersController) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	p := c.useCase.Handle(usecase.GetProductQuery{2})

	switch {
	case p != nil:
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(p)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}
