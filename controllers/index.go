package controllers

import (
	"net/http"

	ctx "github.com/jdkanani/smalldocs/context"
)

// Index handler
func Index(context *ctx.Context, w http.ResponseWriter, r *http.Request) (code int, err error) {
	if r.URL.Path != "/" {
		return http.StatusNotFound, nil
	}
	return 200, context.RenderTemplate(w, "index", nil)
}
