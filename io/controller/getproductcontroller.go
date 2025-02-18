package controller

import (
	"exampleapp/domain/errors"
	"exampleapp/domain/usecase"
	"net/http"
	"strconv"
)

type GetProductController struct {
	errFactory errors.Factory
	useCase    *usecase.GetProductUseCase
}

func NewGetProductController(errFactory errors.Factory, useCase *usecase.GetProductUseCase) *GetProductController {
	return &GetProductController{errFactory, useCase}
}

func (c *GetProductController) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	id, err := strconv.ParseUint(req.PathValue("id"), 10, 64)

	if err != nil {
		HttpJsonError(w, c.errFactory.WrapContext(req.Context(), "GetProduct: %w", err).Error(), http.StatusBadRequest)
		return
	}

	p := c.useCase.Handle(req.Context(), usecase.GetProductQuery{id})

	if p == nil {
		HttpJsonError(w, "Product not found", http.StatusNotFound)
		return
	}

	HttpJson(w, p, http.StatusOK)
}
