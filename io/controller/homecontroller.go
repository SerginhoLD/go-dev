package controller

import (
	"encoding/json"
	"exampleapp/domain/usecase"
	"net/http"
	"strconv"
)

type HomeController struct {
	useCase *usecase.PaginateProductsUseCase
}

func NewHomeController(useCase *usecase.PaginateProductsUseCase) *HomeController {
	return &HomeController{useCase}
}

func (c *HomeController) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	page, err := strconv.ParseUint(req.URL.Query().Get("page"), 10, 64)

	switch {
	case err != nil:
		page = 1
	case page == 0:
		page = 1
	}

	models := c.useCase.Handle(usecase.PaginateProductsQuery{Page: page, Limit: 2})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models)
}
