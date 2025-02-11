package controller

import (
	"encoding/json"
	"exampleapp/domain/usecase"
	"net/http"
)

type HomeController struct {
	useCase *usecase.AllProductsUseCase
}

func NewHomeController(useCase *usecase.AllProductsUseCase) *HomeController {
	return &HomeController{useCase}
}

func (c *HomeController) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	models := c.useCase.Handle()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models)
}
