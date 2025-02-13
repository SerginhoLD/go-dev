package controller

import (
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
	page, _ := strconv.ParseUint(req.URL.Query().Get("page"), 10, 64)

	if page == 0 {
		page = 1
	}

	models := c.useCase.Handle(usecase.PaginateProductsQuery{Page: page, Limit: 2})

	HttpJson(w, models, http.StatusOK)
}
