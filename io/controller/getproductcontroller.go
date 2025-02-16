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
		// > example test log error
		err2 := c.errFactory.NewContext(req.Context(), err.Error())
		// < example test log error

		HttpJsonError(w, err2.Error(), http.StatusBadRequest)
		return
	}

	p := c.useCase.Handle(req.Context(), usecase.GetProductQuery{id})

	if p == nil {
		HttpJsonError(w, "Product not found", http.StatusNotFound)
		return
	}

	HttpJson(w, p, http.StatusOK)
}
