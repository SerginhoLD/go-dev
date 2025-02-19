package app

import (
	"fmt"
	"net/http"
	"os"
)

type CoverageController struct {
}

func NewCoverageController() *CoverageController {
	return &CoverageController{}
}

func (c *CoverageController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	file, err := os.ReadFile("coverage.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	fmt.Fprintln(w, string(file))
}
