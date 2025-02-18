package controller

import (
	"encoding/json"
	"exampleapp/domain/errors"
	"exampleapp/domain/usecase"
	"net/http"
)

type CreateProductController struct {
	errFactory errors.Factory
	useCase    *usecase.CreateProductUseCase
}

func NewCreateProductController(errFactory errors.Factory, useCase *usecase.CreateProductUseCase) *CreateProductController {
	return &CreateProductController{errFactory, useCase}
}

func (c *CreateProductController) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var command usecase.CreateProductCommand
	err := json.NewDecoder(req.Body).Decode(&command)

	if err != nil {
		HttpJsonError(w, c.errFactory.WrapContext(req.Context(), "CreateProduct: %w", err).Error(), http.StatusBadRequest)
		return
	}

	result, err := c.useCase.Handle(req.Context(), command)

	if err != nil {
		HttpJsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	HttpJson(w, result, http.StatusOK)
}
