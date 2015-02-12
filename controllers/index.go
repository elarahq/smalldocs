package controllers

import (
	"net/http"

	ctx "github.com/jdkanani/smalldocs/context"
)

// Index handler
func Index(context *ctx.Context, w http.ResponseWriter, r *http.Request) (code int, err error) {
	return 200, context.RenderTemplate(w, "index", "")
}

// Not found handler
func NotFound(context *ctx.Context, w http.ResponseWriter, r *http.Request) (code int, err error) {
	return 404, nil
}
