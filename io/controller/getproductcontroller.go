package controller

import (
	"encoding/json"
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
	w.Header().Set("Content-Type", "application/json")
	//w.Header().Set("X-Abc", req.PathValue("id"))

	id, err := strconv.ParseUint(req.PathValue("id"), 10, 64)

	if err != nil {
		HttpJsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	p := c.useCase.Handle(usecase.GetProductQuery{id})

	switch {
	case p != nil:
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(p)
	default:
		HttpJsonError(w, "Product not found", http.StatusNotFound)
	}
}
