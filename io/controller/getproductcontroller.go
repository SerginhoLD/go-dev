package controller

import (
	"exampleapp/domain/usecase"
	"net/http"
	"strconv"
)

type GetProductController struct {
	useCase *usecase.GetProductUseCase
}

func NewGetProductController(useCase *usecase.GetProductUseCase) *GetProductController {
	return &GetProductController{useCase}
}

func (c *GetProductController) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	id, err := strconv.ParseUint(req.PathValue("id"), 10, 64)

	if err != nil {
		HttpJsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	p := c.useCase.Handle(usecase.GetProductQuery{id})

	if p == nil {
		HttpJsonError(w, "Product not found", http.StatusNotFound)
		return
	}

	HttpJson(w, p, http.StatusOK)
}
