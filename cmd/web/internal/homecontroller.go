package internal

import (
	"exampleapp/internal/domain/repository"
	"exampleapp/internal/domain/usecase"
	"net/http"

	"github.com/gorilla/schema"
)

type HomeController struct {
	useCase *usecase.PaginateObjectsUseCase
	decoder *schema.Decoder
}

func NewHomeController(useCase *usecase.PaginateObjectsUseCase, decoder *schema.Decoder) *HomeController {
	return &HomeController{useCase, decoder}
}

func (c *HomeController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	query := repository.ObjectsQuery{Limit: 10}
	err := c.decoder.Decode(&query, r.URL.Query())

	if err != nil {
		HttpJsonError(w, err.Error(), 400)
		return
	}

	if query.Page < 1 {
		query.Page = 1
	}

	HttpHtml(w, r, "web/templates/home/index.gohtml", map[string]any{"Query": query, "Data": c.useCase.Handle(r.Context(), query)}, 200)
}
