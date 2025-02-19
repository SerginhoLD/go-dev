package app

import (
	"exampleapp/internal/domain/usecase"
	"net/http"
	"strconv"
)

type HomeController struct {
	useCase *usecase.PaginateProductsUseCase
}

func NewHomeController(useCase *usecase.PaginateProductsUseCase) *HomeController {
	return &HomeController{useCase}
}

func (c *HomeController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.ParseUint(r.URL.Query().Get("page"), 10, 64)

	if page == 0 {
		page = 1
	}

	models := c.useCase.Handle(r.Context(), usecase.PaginateProductsQuery{Page: page, Limit: 2})

	HttpJson(w, models, http.StatusOK)
}
