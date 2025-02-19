package app

import (
	"exampleapp/internal/domain/errors"
	"exampleapp/internal/domain/usecase"
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

func (c *GetProductController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)

	if err != nil {
		HttpJsonError(w, c.errFactory.WrapContext(r.Context(), "GetProduct: %w", err).Error(), http.StatusBadRequest)
		return
	}

	p := c.useCase.Handle(r.Context(), usecase.GetProductQuery{id})

	if p == nil {
		HttpJsonError(w, "Product not found", http.StatusNotFound)
		return
	}

	HttpJson(w, p, http.StatusOK)
}
