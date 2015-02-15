package controllers

import (
	"net/http"

	"github.com/jdkanani/goa"
)

//
// Document index
//
func DocumentIndex(ctx *goa.Context, w http.ResponseWriter, r *http.Request) (int, error) {
	project, err := ProjectByName(ctx.Params["pname"])
	if err != nil {
		return 404, err
	}
	return 200, ctx.Render("docs", map[string]interface{}{
		"project": project,
	})
}

//
//  Page index
//
func PageIndex(ctx *goa.Context, w http.ResponseWriter, r *http.Request) (int, error) {
	project, err := ProjectByName(ctx.Params["pname"])
	if err != nil {
		return 404, err
	}
	return 200, ctx.Render("docs", map[string]interface{}{
		"project": project,
	})
}
