package app

import (
	"encoding/json"
	"exampleapp/internal/domain/errors"
	"exampleapp/internal/domain/usecase"
	"net/http"
)

type CreateProductController struct {
	errFactory errors.Factory
	useCase    *usecase.CreateProductUseCase
}

func NewCreateProductController(errFactory errors.Factory, useCase *usecase.CreateProductUseCase) *CreateProductController {
	return &CreateProductController{errFactory, useCase}
}

func (c *CreateProductController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var command usecase.CreateProductCommand
	err := json.NewDecoder(r.Body).Decode(&command)

	if err != nil {
		HttpJsonError(w, c.errFactory.WrapContext(r.Context(), "CreateProduct: %w", err).Error(), http.StatusBadRequest)
		return
	}

	result, err := c.useCase.Handle(r.Context(), command)

	if err != nil {
		HttpJsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	HttpJson(w, result, http.StatusOK)
}
