package controllers

import (
	"net/http"

	"github.com/loconsolutions/goa"
)

// Slug regxp string
const SLUG = `[a-zA-Z0-9\.-]+`
const ID = `[a-fA-F\d]{24}`

// Index handler
func Index(context *goa.Context, w http.ResponseWriter, r *http.Request) (code int, err error) {
	return 200, context.Render("index", "")
}

// Not found handler
func NotFound(context *goa.Context, w http.ResponseWriter, r *http.Request) (code int, err error) {
	return 404, nil
}
